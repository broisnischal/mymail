CREATE TABLE IF NOT EXISTS "email_metadata" (
	"id" text PRIMARY KEY NOT NULL,
	"email_id" text NOT NULL,
	"headers" jsonb NOT NULL,
	"attachments" jsonb,
	CONSTRAINT "email_metadata_email_id_unique" UNIQUE("email_id")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "emails" (
	"id" text PRIMARY KEY NOT NULL,
	"mailbox_id" text NOT NULL,
	"message_id" varchar(512) NOT NULL,
	"from" varchar(255) NOT NULL,
	"to" jsonb NOT NULL,
	"cc" jsonb,
	"bcc" jsonb,
	"subject" text,
	"text_body" text,
	"html_body" text,
	"minio_path" text NOT NULL,
	"size" integer NOT NULL,
	"received_at" timestamp DEFAULT now() NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "mailboxes" (
	"id" text PRIMARY KEY NOT NULL,
	"user_id" text NOT NULL,
	"address" varchar(255) NOT NULL,
	"is_alias" boolean DEFAULT false NOT NULL,
	"is_temp" boolean DEFAULT false NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "queue_jobs" (
	"id" text PRIMARY KEY NOT NULL,
	"type" varchar(50) NOT NULL,
	"payload" jsonb NOT NULL,
	"status" varchar(20) DEFAULT 'pending' NOT NULL,
	"attempts" integer DEFAULT 0 NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"processed_at" timestamp
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "users" (
	"id" text PRIMARY KEY NOT NULL,
	"email" varchar(255) NOT NULL,
	"password_hash" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "users_email_unique" UNIQUE("email")
);
--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "email_metadata_email_id_idx" ON "email_metadata" ("email_id");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "emails_mailbox_id_idx" ON "emails" ("mailbox_id");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "emails_message_id_idx" ON "emails" ("message_id");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "emails_received_at_idx" ON "emails" ("received_at");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "mailboxes_address_idx" ON "mailboxes" ("address");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "mailboxes_user_id_idx" ON "mailboxes" ("user_id");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "queue_jobs_status_idx" ON "queue_jobs" ("status");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "queue_jobs_type_idx" ON "queue_jobs" ("type");--> statement-breakpoint
CREATE INDEX IF NOT EXISTS "users_email_idx" ON "users" ("email");--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "email_metadata" ADD CONSTRAINT "email_metadata_email_id_emails_id_fk" FOREIGN KEY ("email_id") REFERENCES "emails"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "emails" ADD CONSTRAINT "emails_mailbox_id_mailboxes_id_fk" FOREIGN KEY ("mailbox_id") REFERENCES "mailboxes"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "mailboxes" ADD CONSTRAINT "mailboxes_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE cascade ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
