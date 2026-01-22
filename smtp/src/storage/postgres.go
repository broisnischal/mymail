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

func (p *Postgres) FindMailbox(address string) (*Mailbox, error) {
	var mailbox Mailbox
	query := `SELECT id, user_id, address, is_alias, is_temp, created_at, updated_at 
	          FROM mailboxes WHERE address = $1 LIMIT 1`

	err := p.db.Get(&mailbox, query, address)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &mailbox, nil
}

func (p *Postgres) CreateTempMailbox(address string) (*Mailbox, error) {
	var mailbox Mailbox
	query := `INSERT INTO mailboxes (id, user_id, address, is_alias, is_temp, created_at, updated_at)
	          VALUES (gen_random_uuid(), gen_random_uuid(), $1, false, true, NOW(), NOW())
	          RETURNING id, user_id, address, is_alias, is_temp, created_at, updated_at`

	err := p.db.Get(&mailbox, query, address)
	if err != nil {
		return nil, err
	}
	return &mailbox, nil
}

func (p *Postgres) CreateQueueJob(jobType string, payload map[string]interface{}) error {
	// Marshal payload to JSON for JSONB column
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := `INSERT INTO queue_jobs (id, type, payload, status, attempts, created_at)
	          VALUES (gen_random_uuid(), $1, $2::jsonb, 'pending', 0, NOW())`

	_, err = p.db.Exec(query, jobType, payloadJSON)
	return err
}

type Mailbox struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Address   string    `db:"address"`
	IsAlias   bool      `db:"is_alias"`
	IsTemp    bool      `db:"is_temp"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
