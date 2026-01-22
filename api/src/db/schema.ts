import { pgTable, text, timestamp, integer, boolean, jsonb, varchar, index } from 'drizzle-orm/pg-core';
import { relations } from 'drizzle-orm';

export const users = pgTable('users', {
  id: text('id').primaryKey().$defaultFn(() => crypto.randomUUID()),
  email: varchar('email', { length: 255 }).notNull().unique(),
  passwordHash: text('password_hash').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
}, (table) => ({
  emailIdx: index('users_email_idx').on(table.email),
}));

export const mailboxes = pgTable('mailboxes', {
  id: text('id').primaryKey().$defaultFn(() => crypto.randomUUID()),
  userId: text('user_id').references(() => users.id, { onDelete: 'cascade' }).notNull(),
  address: varchar('address', { length: 255 }).notNull(),
  isAlias: boolean('is_alias').default(false).notNull(),
  isTemp: boolean('is_temp').default(false).notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
}, (table) => ({
  addressIdx: index('mailboxes_address_idx').on(table.address),
  userIdIdx: index('mailboxes_user_id_idx').on(table.userId),
}));

export const emails = pgTable('emails', {
  id: text('id').primaryKey().$defaultFn(() => crypto.randomUUID()),
  mailboxId: text('mailbox_id').references(() => mailboxes.id, { onDelete: 'cascade' }).notNull(),
  messageId: varchar('message_id', { length: 512 }).notNull(),
  from: varchar('from', { length: 255 }).notNull(),
  to: jsonb('to').$type<string[]>().notNull(),
  cc: jsonb('cc').$type<string[]>(),
  bcc: jsonb('bcc').$type<string[]>(),
  subject: text('subject'),
  textBody: text('text_body'),
  htmlBody: text('html_body'),
  minioPath: text('minio_path').notNull(),
  size: integer('size').notNull(),
  receivedAt: timestamp('received_at').defaultNow().notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
}, (table) => ({
  mailboxIdIdx: index('emails_mailbox_id_idx').on(table.mailboxId),
  messageIdIdx: index('emails_message_id_idx').on(table.messageId),
  receivedAtIdx: index('emails_received_at_idx').on(table.receivedAt),
}));

export const emailMetadata = pgTable('email_metadata', {
  id: text('id').primaryKey().$defaultFn(() => crypto.randomUUID()),
  emailId: text('email_id').references(() => emails.id, { onDelete: 'cascade' }).notNull().unique(),
  headers: jsonb('headers').$type<Record<string, string>>().notNull(),
  attachments: jsonb('attachments').$type<Array<{
    filename: string;
    contentType: string;
    size: number;
    minioPath: string;
  }>>(),
}, (table) => ({
  emailIdIdx: index('email_metadata_email_id_idx').on(table.emailId),
}));

export const queueJobs = pgTable('queue_jobs', {
  id: text('id').primaryKey().$defaultFn(() => crypto.randomUUID()),
  type: varchar('type', { length: 50 }).notNull(),
  payload: jsonb('payload').notNull(),
  status: varchar('status', { length: 20 }).default('pending').notNull(),
  attempts: integer('attempts').default(0).notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  processedAt: timestamp('processed_at'),
}, (table) => ({
  statusIdx: index('queue_jobs_status_idx').on(table.status),
  typeIdx: index('queue_jobs_type_idx').on(table.type),
}));

// Relations
export const usersRelations = relations(users, ({ many }) => ({
  mailboxes: many(mailboxes),
}));

export const mailboxesRelations = relations(mailboxes, ({ one, many }) => ({
  user: one(users, {
    fields: [mailboxes.userId],
    references: [users.id],
  }),
  emails: many(emails),
}));

export const emailsRelations = relations(emails, ({ one }) => ({
  mailbox: one(mailboxes, {
    fields: [emails.mailboxId],
    references: [mailboxes.id],
  }),
  metadata: one(emailMetadata, {
    fields: [emails.id],
    references: [emailMetadata.emailId],
  }),
}));
