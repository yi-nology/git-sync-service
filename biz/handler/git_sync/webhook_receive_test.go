package git_sync

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := newRateLimiter(5) // 5 requests per second

	// Should allow first 5 requests immediately
	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	// 6th request should be denied (tokens exhausted)
	if rl.Allow() {
		t.Fatal("expected 6th request to be denied")
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	rl := newRateLimiter(10) // 10 requests per second

	// Exhaust all tokens
	for i := 0; i < 10; i++ {
		rl.Allow()
	}

	// Should be denied now
	if rl.Allow() {
		t.Fatal("expected request to be denied after exhausting tokens")
	}

	// Manually simulate time passage by manipulating lastRefill
	rl.mu.Lock()
	rl.lastRefill = rl.lastRefill.Add(-1 * time.Second) // simulate 1 second passing
	rl.mu.Unlock()

	// Should be allowed now (refilled 10 tokens)
	if !rl.Allow() {
		t.Fatal("expected request to be allowed after refill")
	}
}

func TestRateLimiter_DefaultRate(t *testing.T) {
	rl := newRateLimiter(1) // 1 request per second

	// First request should be allowed
	if !rl.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	// Second immediate request should be denied
	if rl.Allow() {
		t.Fatal("expected second immediate request to be denied")
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := newRateLimiter(100) // 100 requests per second

	var wg sync.WaitGroup
	allowed := make(chan bool, 200)

	// Launch 200 concurrent requests
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed <- rl.Allow()
		}()
	}

	wg.Wait()
	close(allowed)

	allowedCount := 0
	for a := range allowed {
		if a {
			allowedCount++
		}
	}

	// Should allow exactly 100 requests (the rate limit)
	if allowedCount != 100 {
		t.Fatalf("expected exactly 100 allowed requests, got %d", allowedCount)
	}
}

func TestRateLimitMiddleware_AllowsWithinLimit(t *testing.T) {
	// Set up a rate limiter with a high limit for testing
	setWebhookRateLimiter(newRateLimiter(10))
	defer resetWebhookRateLimiter()

	// Create a test context
	ctx := app.NewContext(0)

	nextCalled := false
	nextHandler := func(c context.Context, ctx *app.RequestContext) {
		nextCalled = true
	}

	// Create the middleware
	middleware := RateLimitMiddleware()

	// Should allow the request
	middleware(context.Background(), ctx)

	// Call next handler if not aborted
	if !ctx.IsAborted() {
		nextHandler(context.Background(), ctx)
	}

	if !nextCalled {
		t.Error("expected next handler to be called within rate limit")
	}

	if ctx.IsAborted() {
		t.Error("expected request to not be aborted within rate limit")
	}
}

func TestRateLimitMiddleware_DeniesOverLimit(t *testing.T) {
	// Set up a rate limiter with limit of 1
	setWebhookRateLimiter(newRateLimiter(1))
	defer resetWebhookRateLimiter()

	middleware := RateLimitMiddleware()

	// First request should be allowed
	ctx1 := app.NewContext(0)
	middleware(context.Background(), ctx1)
	if ctx1.IsAborted() {
		t.Fatal("expected first request to be allowed")
	}

	// Second immediate request should be denied
	ctx2 := app.NewContext(0)
	middleware(context.Background(), ctx2)
	if !ctx2.IsAborted() {
		t.Fatal("expected second request to be aborted (rate limited)")
	}

	// Verify 429 status code
	if ctx2.Response.StatusCode() != consts.StatusTooManyRequests {
		t.Errorf("expected status code %d, got %d", consts.StatusTooManyRequests, ctx2.Response.StatusCode())
	}

	// Verify error message
	body := string(ctx2.Response.Body())
	if body == "" {
		t.Error("expected error message in response body")
	}
}

func TestRateLimitMiddleware_AbortPreventsNext(t *testing.T) {
	// Set up a rate limiter with limit of 1 and exhaust it
	rl := newRateLimiter(1)
	rl.Allow() // exhaust the token
	setWebhookRateLimiter(rl)
	defer resetWebhookRateLimiter()

	middleware := RateLimitMiddleware()

	ctx := app.NewContext(0)
	nextCalled := false
	nextHandler := func(c context.Context, ctx *app.RequestContext) {
		nextCalled = true
	}

	middleware(context.Background(), ctx)

	// Call next handler if not aborted
	if !ctx.IsAborted() {
		nextHandler(context.Background(), ctx)
	}

	if nextCalled {
		t.Error("expected next handler to NOT be called when rate limited")
	}

	if !ctx.IsAborted() {
		t.Error("expected request to be aborted when rate limited")
	}
}

func TestRateLimitMiddleware_ReturnsCorrectStatusCode(t *testing.T) {
	// Set up a rate limiter with limit of 1 and exhaust it
	rl := newRateLimiter(1)
	rl.Allow() // exhaust the token
	setWebhookRateLimiter(rl)
	defer resetWebhookRateLimiter()

	middleware := RateLimitMiddleware()
	ctx := app.NewContext(0)
	middleware(context.Background(), ctx)

	if ctx.Response.StatusCode() != http.StatusTooManyRequests {
		t.Errorf("expected status code %d (TooManyRequests), got %d", http.StatusTooManyRequests, ctx.Response.StatusCode())
	}
}

func TestRateLimiter_ZeroRate(t *testing.T) {
	// A rate limiter with 0 rate should never allow
	rl := newRateLimiter(0)

	if rl.Allow() {
		t.Fatal("expected request to be denied with zero rate")
	}
}

func TestRateLimiter_LargeBurst(t *testing.T) {
	rl := newRateLimiter(1000) // 1000 requests per second

	// Should allow 1000 requests
	allowed := 0
	for i := 0; i < 1000; i++ {
		if rl.Allow() {
			allowed++
		}
	}

	if allowed != 1000 {
		t.Fatalf("expected 1000 allowed requests, got %d", allowed)
	}

	// Note: Due to token refill based on elapsed time, the 1001st request
	// might be allowed if enough time has passed. This is expected behavior
	// for a token bucket rate limiter.
}

func TestRateLimiter_ConcurrentSafety(t *testing.T) {
	rl := newRateLimiter(50)

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]bool, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result := rl.Allow()
			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	allowedCount := 0
	for _, r := range results {
		if r {
			allowedCount++
		}
	}

	// Should allow exactly 50 requests
	if allowedCount != 50 {
		t.Fatalf("expected exactly 50 allowed requests, got %d", allowedCount)
	}
}
