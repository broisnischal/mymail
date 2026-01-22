import { Hono } from 'hono';
import { z } from 'zod';
import { db } from '../db';
import { mailboxes } from '../db/schema';
import { eq, and } from 'drizzle-orm';
import { authMiddleware } from '../middleware/auth';
import { config } from '@shared/config';

const app = new Hono();

app.use('/*', authMiddleware);

const createMailboxSchema = z.object({
  address: z.string().email(),
  isAlias: z.boolean().optional().default(false),
});

app.get('/', async (c) => {
  const userId = c.get('userId');
  const userMailboxes = await db.select().from(mailboxes).where(eq(mailboxes.userId, userId));
  return c.json({ mailboxes: userMailboxes });
});

app.post('/', async (c) => {
  try {
    const userId = c.get('userId');
    const body = await c.req.json();
    const { address, isAlias } = createMailboxSchema.parse(body);

    // Validate domain
    const domain = address.split('@')[1];
    if (domain !== config.smtp.domain) {
      return c.json({ error: 'Invalid domain' }, 400);
    }

    // Check if mailbox exists
    const existing = await db.select().from(mailboxes).where(eq(mailboxes.address, address)).limit(1);
    if (existing.length > 0) {
      return c.json({ error: 'Mailbox already exists' }, 400);
    }

    const [mailbox] = await db.insert(mailboxes).values({
      userId,
      address,
      isAlias,
      isTemp: false,
    }).returning();

    return c.json({ mailbox });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return c.json({ error: error.errors }, 400);
    }
    return c.json({ error: 'src server error' }, 500);
  }
});

app.delete('/:id', async (c) => {
  const userId = c.get('userId');
  const id = c.req.param('id');

  const [mailbox] = await db.select().from(mailboxes)
    .where(and(eq(mailboxes.id, id), eq(mailboxes.userId, userId)))
    .limit(1);

  if (!mailbox) {
    return c.json({ error: 'Mailbox not found' }, 404);
  }

  await db.delete(mailboxes).where(eq(mailboxes.id, id));
  return c.json({ success: true });
});

export default app;
