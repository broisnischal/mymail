import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { serve } from '@hono/node-server';
import { config } from '@shared/config';
import authRoutes from './routes/auth';
import mailboxRoutes from './routes/mailboxes';
import emailRoutes from './routes/emails';
import { ensureBucket } from './services/minio';

const app = new Hono();

app.use('/*', cors({
  origin: '*',
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowHeaders: ['Content-Type', 'Authorization'],
}));

app.get('/health', (c) => {
  return c.json({ status: 'ok', timestamp: new Date().toISOString() });
});

app.route('/api/auth', authRoutes);
app.route('/api/mailboxes', mailboxRoutes);
app.route('/api/emails', emailRoutes);

// Initialize MinIO bucket
ensureBucket().catch(console.error);

const port = config.api.port;
console.log(`🚀 API Server running on http://${config.api.host}:${port}`);

serve({
  fetch: app.fetch,
  port,
});
