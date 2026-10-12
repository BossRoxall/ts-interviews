package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

// fakeClock only moves when the test calls Advance.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// 3 calls per minute, refilling 1 call every 20 seconds.
var testConfig = Config{
	Capacity:       3,
	Interval:       time.Minute,
	RefillAmount:   1,
	RefillInterval: 20 * time.Second,
}

// Several users so that state leaking between keys shows up as a failure.
var users = []string{"alice", "bob", "carol"}

func TestLimiterDrainLimitAndRefillPerUser(t *testing.T) {
	clock := newFakeClock()
	limiter, err := New(testConfig, clock)
	if err != nil {
		t.Fatalf("New(%+v): unexpected error: %v", testConfig, err)
	}

	t.Run("1 each user can spend their full capacity", func(t *testing.T) {
		for _, user := range users {
			for call := 1; call <= testConfig.Capacity; call++ {
				if !limiter.Allow(user) {
					t.Errorf("%s: call %d of %d was denied, want allowed", user, call, testConfig.Capacity)
				}
			}
		}
	})

	t.Run("2 each user is denied once over the limit", func(t *testing.T) {
		for _, user := range users {
			if limiter.Allow(user) {
				t.Errorf("%s: call %d was allowed, want denied (capacity is %d)", user, testConfig.Capacity+1, testConfig.Capacity)
			}
		}
	})

	clock.Advance(testConfig.RefillInterval)

	t.Run("3 each user gets one call back after one refill interval", func(t *testing.T) {
		for _, user := range users {
			if !limiter.Allow(user) {
				t.Errorf("%s: first call after %s was denied, want allowed", user, testConfig.RefillInterval)
			}
		}
	})

	t.Run("4 each user is denied again after spending the refill", func(t *testing.T) {
		for _, user := range users {
			if limiter.Allow(user) {
				t.Errorf("%s: second call after %s was allowed, want denied", user, testConfig.RefillInterval)
			}
		}
	})
}
