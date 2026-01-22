package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/mymail/smtp/src/config"
	"github.com/mymail/smtp/src/ratelimit"
	"github.com/mymail/smtp/src/storage"
)

type Backend struct {
	db          *storage.Postgres
	redis       *storage.Redis
	minio       *storage.MinIO
	rateLimiter *ratelimit.RateLimiter
	cfg         *config.Config
}

func NewBackend(db *storage.Postgres, redis *storage.Redis, minio *storage.MinIO,
	rateLimiter *ratelimit.RateLimiter, cfg *config.Config) *Backend {
	return &Backend{
		db:          db,
		redis:       redis,
		minio:       minio,
		rateLimiter: rateLimiter,
		cfg:         cfg,
	}
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	remoteAddr := c.Conn().RemoteAddr().String()
	ip := strings.Split(remoteAddr, ":")[0]

	// Rate limit by IP
	ctx := context.Background()
	allowed, err := b.rateLimiter.AllowConnection(ctx, ip)
	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded")
	}
	if err != nil {
		return nil, err
	}

	return &Session{
		backend:    b,
		remoteAddr: ip,
	}, nil
}

type Session struct {
	backend    *Backend
	remoteAddr string
	from       string
	to         []string
}

func (s *Session) AuthMechanism() []string {
	return []string{}
}

func (s *Session) AuthPlain(username, password string) error {
	// No authentication required for receiving emails
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// Extract domain from recipient
	parts := strings.Split(to, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email address")
	}

	domain := parts[1]
	if domain != s.backend.cfg.SMTP.Domain {
		return fmt.Errorf("invalid domain")
	}

	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	ctx := context.Background()

	// Use io.TeeReader to split the stream: one for header parsing, one for MinIO streaming
	// This allows us to read headers while simultaneously streaming to MinIO
	headerBuf := &bytes.Buffer{}
	teeReader := io.TeeReader(r, headerBuf)

	headerLimit := int64(64 * 1024)
	headerReader := io.LimitReader(teeReader, headerLimit)

	headerSize, err := io.Copy(headerBuf, headerReader)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	msg, err := message.Read(bytes.NewReader(headerBuf.Bytes()))
	if err != nil {
		return fmt.Errorf("failed to parse message headers: %w", err)
	}

	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")
	messageID := msg.Header.Get("Message-ID")
	if messageID == "" {
		messageID = fmt.Sprintf("<%s@%s>", uuid.New().String(), s.backend.cfg.SMTP.Domain)
	}

	toAddresses := []string{}
	if to != "" {
		addresses, _ := mail.ParseAddressList(to)
		for _, addr := range addresses {
			toAddresses = append(toAddresses, addr.Address)
		}
	}

	validMailboxes := []*storage.Mailbox{}
	for _, recipient := range s.to {
		mailbox, err := s.backend.db.FindMailbox(recipient)
		if err != nil || mailbox == nil {
			continue
		}

		// Rate limit check
		allowed, err := s.backend.rateLimiter.AllowEmail(ctx, mailbox.UserID)
		if !allowed || err != nil {
			continue
		}

		validMailboxes = append(validMailboxes, mailbox)
	}

	if len(validMailboxes) == 0 {
		return nil // No valid recipients
	}

	// For streaming: combine header buffer with remaining stream
	// The teeReader continues reading from 'r' after headers
	var fullStream io.Reader
	if headerSize < headerLimit {
		// All data was in header buffer (small email)
		fullStream = bytes.NewReader(headerBuf.Bytes())
	} else {
		// We have more data - combine header buffer with remaining from teeReader
		// teeReader will continue reading the rest of the stream
		fullStream = io.MultiReader(bytes.NewReader(headerBuf.Bytes()), teeReader)
	}

	// Upload once to a shared location (use first mailbox's path as primary)
	// For multiple recipients, we'll reference this file
	primaryMailbox := validMailboxes[0]
	emailID := uuid.New().String()
	primaryPath := fmt.Sprintf("%s/%s/%s.eml", primaryMailbox.UserID, time.Now().Format("2006/01/02"), emailID)

	// Stream directly to MinIO - this streams the entire email without buffering
	err = s.backend.minio.UploadStream(ctx, primaryPath, fullStream)
	if err != nil {
		return fmt.Errorf("failed to upload email: %w", err)
	}

	// Extract limited text/HTML bodies from header buffer
	var textBody, htmlBody string
	bodyReader := bytes.NewReader(headerBuf.Bytes())
	if bodyMsg, err := message.Read(bodyReader); err == nil {
		if mr := bodyMsg.MultipartReader(); mr != nil {
			if p, err := mr.NextPart(); err == nil {
				mediaType, _, _ := p.Header.ContentType()
				if body, err := io.ReadAll(io.LimitReader(p.Body, 10240)); err == nil {
					if strings.HasPrefix(mediaType, "text/plain") {
						textBody = string(body)
					} else if strings.HasPrefix(mediaType, "text/html") {
						htmlBody = string(body)
					}
				}
			}
		} else {
			if body, err := io.ReadAll(io.LimitReader(bodyMsg.Body, 10240)); err == nil {
				textBody = string(body)
			}
		}
	}

	// Create queue jobs for all recipients
	// For now, all recipients point to the same file (could be optimized with copies)
	for _, mailbox := range validMailboxes {
		var path string
		if mailbox == primaryMailbox {
			path = primaryPath
		} else {
			// For other recipients, use the same file (shared storage)
			// In production, you might want to copy the file for each recipient
			path = primaryPath
		}

		// Approximate size (will be updated by worker when processing)
		emailSize := int64(headerBuf.Len())

		// Create queue job for processing
		payload := map[string]interface{}{
			"email_id":   emailID,
			"mailbox_id": mailbox.ID,
			"message_id": messageID,
			"from":       from,
			"to":         toAddresses,
			"subject":    subject,
			"text_body":  textBody,
			"html_body":  htmlBody,
			"minio_path": path,
			"size":       emailSize,
		}

		err = s.backend.db.CreateQueueJob("process_email", payload)
		if err != nil {
			continue
		}
	}

	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = nil
}

func (s *Session) Logout() error {
	return nil
}
