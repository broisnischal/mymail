import Redis from 'ioredis';
import { config } from '@shared/config';

export const redis = new Redis(config.redis.url, {
  keyPrefix: config.redis.prefix,
  maxRetriesPerRequest: 3,
});

export class CacheService {
  async get<T>(key: string): Promise<T | null> {
    const value = await redis.get(key);
    return value ? JSON.parse(value) : null;
  }

  async set(key: string, value: any, ttl?: number): Promise<void> {
    const serialized = JSON.stringify(value);
    if (ttl) {
      await redis.setex(key, ttl, serialized);
    } else {
      await redis.set(key, serialized);
    }
  }

  async del(key: string): Promise<void> {
    await redis.del(key);
  }

  async exists(key: string): Promise<boolean> {
    return (await redis.exists(key)) === 1;
  }

  async incr(key: string): Promise<number> {
    return await redis.incr(key);
  }

  async expire(key: string, seconds: number): Promise<void> {
    await redis.expire(key, seconds);
  }
}

export const cache = new CacheService();
