import { Hono } from 'hono';
import { db } from '../db';
import { emails, mailboxes, emailMetadata } from '../db/schema';
import { eq, and, desc } from 'drizzle-orm';
import { authMiddleware } from '../middleware/auth';
import { getEmail } from '../services/minio';

const app = new Hono();

app.use('/*', authMiddleware);

app.get('/', async (c) => {
  const userId = c.get('userId');
  const mailboxId = c.req.query('mailboxId');
  const limit = parseInt(c.req.query('limit') || '50');
  const offset = parseInt(c.req.query('offset') || '0');

  let query = db.select({
    id: emails.id,
    messageId: emails.messageId,
    from: emails.from,
    to: emails.to,
    cc: emails.cc,
    bcc: emails.bcc,
    subject: emails.subject,
    size: emails.size,
    receivedAt: emails.receivedAt,
    createdAt: emails.createdAt,
    mailboxId: emails.mailboxId,
    mailboxAddress: mailboxes.address,
  })
    .from(emails)
    .innerJoin(mailboxes, eq(emails.mailboxId, mailboxes.id))
    .where(eq(mailboxes.userId, userId))
    .orderBy(desc(emails.receivedAt))
    .limit(limit)
    .offset(offset);

  if (mailboxId) {
    query = query.where(and(eq(mailboxes.userId, userId), eq(emails.mailboxId, mailboxId)));
  }

  const emailList = await query;

  return c.json({ emails: emailList });
});

app.get('/:id', async (c) => {
  const userId = c.get('userId');
  const id = c.req.param('id');

  const [email] = await db.select({
    id: emails.id,
    messageId: emails.messageId,
    from: emails.from,
    to: emails.to,
    cc: emails.cc,
    bcc: emails.bcc,
    subject: emails.subject,
    textBody: emails.textBody,
    htmlBody: emails.htmlBody,
    size: emails.size,
    receivedAt: emails.receivedAt,
    createdAt: emails.createdAt,
    minioPath: emails.minioPath,
    mailboxId: emails.mailboxId,
    mailboxAddress: mailboxes.address,
  })
    .from(emails)
    .innerJoin(mailboxes, eq(emails.mailboxId, mailboxes.id))
    .where(and(eq(emails.id, id), eq(mailboxes.userId, userId)))
    .limit(1);

  if (!email) {
    return c.json({ error: 'Email not found' }, 404);
  }

  // Get metadata
  const [metadata] = await db.select().from(emailMetadata)
    .where(eq(emailMetadata.emailId, id))
    .limit(1);

  return c.json({
    email: {
      ...email,
      metadata: metadata || null,
    },
  });
});

app.get('/:id/raw', async (c) => {
  const userId = c.get('userId');
  const id = c.req.param('id');

  const [email] = await db.select({
    minioPath: emails.minioPath,
    mailboxId: emails.mailboxId,
  })
    .from(emails)
    .innerJoin(mailboxes, eq(emails.mailboxId, mailboxes.id))
    .where(and(eq(emails.id, id), eq(mailboxes.userId, userId)))
    .limit(1);

  if (!email) {
    return c.json({ error: 'Email not found' }, 404);
  }

  const rawEmail = await getEmail(email.minioPath);
  c.header('Content-Type', 'message/rfc822');
  return c.body(rawEmail);
});

app.delete('/:id', async (c) => {
  const userId = c.get('userId');
  const id = c.req.param('id');

  const [email] = await db.select({
    id: emails.id,
    minioPath: emails.minioPath,
  })
    .from(emails)
    .innerJoin(mailboxes, eq(emails.mailboxId, mailboxes.id))
    .where(and(eq(emails.id, id), eq(mailboxes.userId, userId)))
    .limit(1);

  if (!email) {
    return c.json({ error: 'Email not found' }, 404);
  }

  await db.delete(emails).where(eq(emails.id, id));
  return c.json({ success: true });
});

export default app;
