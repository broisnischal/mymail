package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/mymail/smtp/internal/storage"
)

type RateLimiter struct {
	redis *storage.Redis
}

func New(redis *storage.Redis) *RateLimiter {
	return &RateLimiter{redis: redis}
}

func (r *RateLimiter) AllowConnection(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("ratelimit:connection:%s", ip)
	count, err := r.redis.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.redis.Expire(ctx, key, time.Minute)
	}

	// Allow up to 10 connections per minute per IP
	return count <= 10, nil
}

func (r *RateLimiter) AllowEmail(ctx context.Context, userID string) (bool, error) {
	// Per user rate limit
	userKey := fmt.Sprintf("ratelimit:email:user:%s", userID)
	userCount, err := r.redis.Incr(ctx, userKey)
	if err != nil {
		return false, err
	}

	if userCount == 1 {
		r.redis.Expire(ctx, userKey, 24*time.Hour)
	}

	// Per hour rate limit
	hourKey := fmt.Sprintf("ratelimit:email:hour:%s:%d", userID, time.Now().Unix()/3600)
	hourCount, err := r.redis.Incr(ctx, hourKey)
	if err != nil {
		return false, err
	}

	if hourCount == 1 {
		r.redis.Expire(ctx, hourKey, time.Hour)
	}

	// Check limits
	if userCount > 1000 {
		return false, nil
	}
	if hourCount > 100 {
		return false, nil
	}

	return true, nil
}
