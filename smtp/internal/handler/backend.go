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
	"github.com/mymail/smtp/internal/config"
	"github.com/mymail/smtp/internal/ratelimit"
	"github.com/mymail/smtp/internal/storage"
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
	// Read entire message
	var buf bytes.Buffer
	size, err := io.Copy(&buf, r)
	if err != nil {
		return err
	}

	if size > s.backend.cfg.SMTP.MaxMessageSize {
		return fmt.Errorf("message too large")
	}

	// Parse message
	msg, err := message.Read(&buf)
	if err != nil {
		return err
	}

	// Extract headers
	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")
	messageID := msg.Header.Get("Message-ID")
	if messageID == "" {
		messageID = fmt.Sprintf("<%s@%s>", uuid.New().String(), s.backend.cfg.SMTP.Domain)
	}

	// Parse recipients
	toAddresses := []string{}
	if to != "" {
		addresses, _ := mail.ParseAddressList(to)
		for _, addr := range addresses {
			toAddresses = append(toAddresses, addr.Address)
		}
	}

	// Process each recipient
	ctx := context.Background()
	for _, recipient := range s.to {
		mailbox, err := s.backend.db.FindMailbox(recipient)
		if err != nil {
			continue
		}

		if mailbox == nil {
			continue
		}

		// Rate limit
		allowed, err := s.backend.rateLimiter.AllowEmail(ctx, mailbox.UserID)
		if !allowed || err != nil {
			continue
		}

		// Generate storage path
		emailID := uuid.New().String()
		path := fmt.Sprintf("%s/%s/%s.eml", mailbox.UserID, time.Now().Format("2006/01/02"), emailID)

		// Upload to MinIO
		bufReader := bytes.NewReader(buf.Bytes())
		err = s.backend.minio.Upload(ctx, path, bufReader, size)
		if err != nil {
			continue
		}

		// Extract text and HTML bodies
		var textBody, htmlBody string
		if mr := msg.MultipartReader(); mr != nil {
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					continue
				}

				mediaType, _, _ := p.Header.ContentType()
				body, _ := io.ReadAll(p.Body)

				if strings.HasPrefix(mediaType, "text/plain") {
					textBody = string(body)
				} else if strings.HasPrefix(mediaType, "text/html") {
					htmlBody = string(body)
				}
			}
		} else {
			body, _ := io.ReadAll(msg.Body)
			textBody = string(body)
		}

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
			"size":       size,
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
