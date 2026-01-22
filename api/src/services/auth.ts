import bcrypt from 'bcrypt';
import { sign, verify } from 'hono/jwt';
import { config } from '@shared/config';

export async function hashPassword(password: string): Promise<string> {
  return bcrypt.hash(password, 10);
}

export async function verifyPassword(password: string, hash: string): Promise<boolean> {
  return bcrypt.compare(password, hash);
}

export async function generateToken(userId: string): Promise<string> {
  const payload = {
    userId,
    exp: Math.floor(Date.now() / 1000) + 7 * 24 * 60 * 60,
  };
  return await sign(payload, config.api.jwtSecret);
}

export async function verifyToken(token: string): Promise<{ userId: string } | null> {
  try {
    const payload = await verify(token, config.api.jwtSecret);
    return payload as { userId: string };
  } catch {
    return null;
  }
}
