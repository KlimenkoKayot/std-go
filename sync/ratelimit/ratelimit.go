//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	limit        chan interface{}
	closeCtx     context.Context
	closeCtxFunc context.CancelFunc
	interval     time.Duration
	maxCount     int
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	limit := make(chan interface{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		if interval == 0 {
			for {
				select {
				case limit <- struct{}{}:
					continue
				case <-ctx.Done():
					return
				}
			}
		} else {
			limit <- struct{}{}
			ticker := time.NewTicker(time.Duration(interval.Nanoseconds() / int64(maxCount)))
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					select {
					case limit <- struct{}{}:
						continue
					default:
						continue
					}
				case <-ctx.Done():
					return
				}
			}
		}
	}(ctx)
	return &Limiter{
		limit,
		ctx,
		cancel,
		interval,
		maxCount,
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.closeCtx.Done():
		return ErrStopped
	case <-ctx.Done():
		return ctx.Err()
	default:
		select {
		case <-l.closeCtx.Done():
			return ErrStopped
		case <-ctx.Done():
			return ctx.Err()
		case <-l.limit:
			return nil
		}
	}
}

func (l *Limiter) Stop() {
	l.closeCtxFunc()
}
