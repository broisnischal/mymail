import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from './schema';
import { config } from '@shared/config';

const client = postgres(config.database.url, {
  max: config.database.maxConnections,
});

export const db = drizzle(client, { schema });

export type Database = typeof db;
