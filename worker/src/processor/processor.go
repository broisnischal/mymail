package processor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mymail/worker/src/config"
	"github.com/mymail/worker/src/storage"
)

type Processor struct {
	db     *storage.Postgres
	redis  *storage.Redis
	config *config.Config
}

func New(db *storage.Postgres, redis *storage.Redis, cfg *config.Config) *Processor {
	return &Processor{
		db:     db,
		redis:  redis,
		config: cfg,
	}
}

func (p *Processor) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *Processor) processBatch(ctx context.Context) {
	jobs, err := p.db.GetPendingJobs(p.config.Worker.BatchSize)
	if err != nil {
		log.Printf("Error fetching jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		return
	}

	log.Printf("Processing %d jobs", len(jobs))

	for _, job := range jobs {
		if err := p.processJob(ctx, job); err != nil {
			log.Printf("Error processing job %s: %v", job.ID, err)
			p.db.IncrementJobAttempts(job.ID)
			if job.Attempts >= 3 {
				p.db.UpdateJobStatus(job.ID, "failed")
			}
		} else {
			p.db.UpdateJobStatus(job.ID, "completed")
		}
	}
}

func (p *Processor) processJob(ctx context.Context, job storage.QueueJob) error {
	switch job.Type {
	case "process_email":
		return p.processEmail(ctx, job)
	default:
		log.Printf("Unknown job type: %s", job.Type)
		return nil
	}
}

func (p *Processor) processEmail(ctx context.Context, job storage.QueueJob) error {
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		return err
	}

	// Extract payload data
	emailID, _ := payload["email_id"].(string)
	if emailID == "" {
		emailID = uuid.New().String()
	}

	mailboxID, _ := payload["mailbox_id"].(string)
	messageID, _ := payload["message_id"].(string)
	from, _ := payload["from"].(string)
	subject, _ := payload["subject"].(string)
	textBody, _ := payload["text_body"].(string)
	htmlBody, _ := payload["html_body"].(string)
	minioPath, _ := payload["minio_path"].(string)
	size, _ := payload["size"].(float64)

	to, _ := payload["to"].([]interface{})
	toAddresses := make([]string, 0, len(to))
	for _, addr := range to {
		if str, ok := addr.(string); ok {
			toAddresses = append(toAddresses, str)
		}
	}

	// Create email record
	email := &storage.Email{
		ID:         emailID,
		MailboxID:  mailboxID,
		MessageID:  messageID,
		From:       from,
		To:         toAddresses,
		Subject:    subject,
		TextBody:   textBody,
		HTMLBody:   htmlBody,
		MinIOPath:  minioPath,
		Size:       int64(size),
		ReceivedAt: time.Now(),
	}

	if err := p.db.CreateEmail(email); err != nil {
		return err
	}

	// Create metadata
	headers := make(map[string]interface{})
	if h, ok := payload["headers"].(map[string]interface{}); ok {
		headers = h
	}

	metadata := &storage.EmailMetadata{
		EmailID:     emailID,
		Headers:     headers,
		Attachments: []interface{}{},
	}

	if err := p.db.CreateEmailMetadata(metadata); err != nil {
		return err
	}

	// Publish notification to Redis
	p.redis.Publish(ctx, "email:received", map[string]interface{}{
		"email_id":   emailID,
		"mailbox_id": mailboxID,
	})

	return nil
}
