package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Postgres{db: db}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) GetDB() *sqlx.DB {
	return p.db
}

func (p *Postgres) GetPendingJobs(limit int) ([]QueueJob, error) {
	var jobs []QueueJob
	query := `SELECT id, type, payload, status, attempts, created_at, processed_at
	          FROM queue_jobs 
	          WHERE status = 'pending' 
	          ORDER BY created_at ASC 
	          LIMIT $1`

	err := p.db.Select(&jobs, query, limit)
	return jobs, err
}

func (p *Postgres) UpdateJobStatus(id string, status string) error {
	query := `UPDATE queue_jobs SET status = $1, processed_at = NOW() WHERE id = $2`
	_, err := p.db.Exec(query, status, id)
	return err
}

func (p *Postgres) IncrementJobAttempts(id string) error {
	query := `UPDATE queue_jobs SET attempts = attempts + 1 WHERE id = $2`
	_, err := p.db.Exec(query, id)
	return err
}

func (p *Postgres) CreateEmail(email *Email) error {
	query := `INSERT INTO emails (id, mailbox_id, message_id, "from", "to", cc, bcc, subject, text_body, html_body, minio_path, size, received_at, created_at)
	          VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7::jsonb, $8, $9, $10, $11, $12, $13, NOW())
	          ON CONFLICT (id) DO NOTHING
	          RETURNING id`

	// Convert []string slices to JSON for JSONB columns
	// Always marshal to ensure valid JSON (empty array [] for nil)
	toJSON, _ := json.Marshal(email.To)
	ccJSON, _ := json.Marshal(email.CC)
	bccJSON, _ := json.Marshal(email.BCC)

	err := p.db.Get(&email.ID, query,
		email.ID, email.MailboxID, email.MessageID, email.From, toJSON, ccJSON, bccJSON,
		email.Subject, email.TextBody, email.HTMLBody, email.MinIOPath, email.Size, email.ReceivedAt)

	// If no rows returned, email already existed (ON CONFLICT DO NOTHING)
	// This is fine - the email was already processed
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (p *Postgres) CreateEmailMetadata(metadata *EmailMetadata) error {
	query := `INSERT INTO email_metadata (id, email_id, headers, attachments, created_at)
	          VALUES (gen_random_uuid(), $1, $2::jsonb, $3::jsonb, NOW())
	          RETURNING id`

	// Marshal to JSON for JSONB columns
	// Ensure we always have valid JSON (empty object {} for nil map, empty array [] for nil slice)
	if metadata.Headers == nil {
		metadata.Headers = make(map[string]interface{})
	}
	if metadata.Attachments == nil {
		metadata.Attachments = []interface{}{}
	}

	headersJSON, _ := json.Marshal(metadata.Headers)
	attachmentsJSON, _ := json.Marshal(metadata.Attachments)

	return p.db.Get(&metadata.ID, query, metadata.EmailID, headersJSON, attachmentsJSON)
}

type QueueJob struct {
	ID          string     `db:"id"`
	Type        string     `db:"type"`
	Payload     string     `db:"payload"`
	Status      string     `db:"status"`
	Attempts    int        `db:"attempts"`
	CreatedAt   time.Time  `db:"created_at"`
	ProcessedAt *time.Time `db:"processed_at"`
}

type Email struct {
	ID         string    `db:"id"`
	MailboxID  string    `db:"mailbox_id"`
	MessageID  string    `db:"message_id"`
	From       string    `db:"from"`
	To         []string  `db:"to"`
	CC         []string  `db:"cc"`
	BCC        []string  `db:"bcc"`
	Subject    string    `db:"subject"`
	TextBody   string    `db:"text_body"`
	HTMLBody   string    `db:"html_body"`
	MinIOPath  string    `db:"minio_path"`
	Size       int64     `db:"size"`
	ReceivedAt time.Time `db:"received_at"`
}

type EmailMetadata struct {
	ID          string                 `db:"id"`
	EmailID     string                 `db:"email_id"`
	Headers     map[string]interface{} `db:"headers"`
	Attachments []interface{}          `db:"attachments"`
	CreatedAt   time.Time              `db:"created_at"`
}
