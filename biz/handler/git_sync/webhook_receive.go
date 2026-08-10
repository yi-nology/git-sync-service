package git_sync

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

const maxWebhookBodySize = 10 << 20

// rateLimiter implements a thread-safe token bucket rate limiter.
type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func newRateLimiter(ratePerSecond int) *rateLimiter {
	return &rateLimiter{
		tokens:     float64(ratePerSecond),
		maxTokens:  float64(ratePerSecond),
		refillRate: float64(ratePerSecond),
		lastRefill: time.Now(),
	}
}

// Allow reports whether a request is permitted under the rate limit.
func (rl *rateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.maxTokens {
		rl.tokens = rl.maxTokens
	}
	rl.lastRefill = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

var (
	webhookRateLimiter   *rateLimiter
	webhookRateLimiterMu sync.RWMutex
	webhookRateLimiterInited bool
)

// getWebhookRateLimiter returns the singleton rate limiter, initializing it from config on first call.
func getWebhookRateLimiter() *rateLimiter {
	webhookRateLimiterMu.RLock()
	if webhookRateLimiterInited {
		defer webhookRateLimiterMu.RUnlock()
		return webhookRateLimiter
	}
	webhookRateLimiterMu.RUnlock()

	webhookRateLimiterMu.Lock()
	defer webhookRateLimiterMu.Unlock()
	// Double-check after acquiring write lock
	if webhookRateLimiterInited {
		return webhookRateLimiter
	}
	rateLimit := GetSyncService().GetConfig().Webhook.RateLimit
	if rateLimit <= 0 {
		rateLimit = 10 // default: 10 requests per second
	}
	webhookRateLimiter = newRateLimiter(rateLimit)
	webhookRateLimiterInited = true
	return webhookRateLimiter
}

// setWebhookRateLimiter overrides the rate limiter. Used for testing.
func setWebhookRateLimiter(rl *rateLimiter) {
	webhookRateLimiterMu.Lock()
	defer webhookRateLimiterMu.Unlock()
	webhookRateLimiter = rl
	webhookRateLimiterInited = true
}

// resetWebhookRateLimiter resets the rate limiter so it will be re-initialized from config. Used for testing.
func resetWebhookRateLimiter() {
	webhookRateLimiterMu.Lock()
	defer webhookRateLimiterMu.Unlock()
	webhookRateLimiter = nil
	webhookRateLimiterInited = false
}

// RateLimitMiddleware returns a middleware that enforces webhook rate limiting.
// When the rate limit is exceeded, it responds with HTTP 429 Too Many Requests.
func RateLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !getWebhookRateLimiter().Allow() {
			response.Error(c, consts.StatusTooManyRequests, "rate limit exceeded, please try again later")
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}

func ReceiveWebhook(ctx context.Context, c *app.RequestContext) {
	repoKey := c.Param("repoKey")
	if repoKey == "" {
		response.BadRequest(c, "repoKey is required")
		return
	}

	bodyBytes, _ := c.Body()
	if len(bodyBytes) > maxWebhookBodySize {
		response.Error(c, consts.StatusRequestEntityTooLarge, "request body too large")
		return
	}

	headers := make(map[string]string)
	c.Request.Header.VisitAll(func(k, v []byte) {
		headers[string(k)] = string(v)
	})

	httpReq, err := http.NewRequest(
		string(c.Method()),
		string(c.Path()),
		io.LimitReader(bytes.NewReader(bodyBytes), maxWebhookBodySize),
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	httpReq.RemoteAddr = c.ClientIP()

	err = GetSyncService().ReceiveWebhook(ctx, repoKey, httpReq)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"message": "webhook received",
	})
}
