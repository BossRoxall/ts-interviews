// Package ratelimiter limits how often each key (for example a user ID) may
// make a call.
package ratelimiter

import (
	"errors"
	"time"

	"ratelimiter/internal/clock"
)

// Config describes the limit applied to every key.
//
// Each key may make up to Capacity calls per Interval, and spent calls come
// back at RefillAmount calls every RefillInterval. For example:
//
//	Config{Capacity: 3, Interval: time.Minute, RefillAmount: 1, RefillInterval: 20 * time.Second}
//
// allows a burst of 3 calls, then one more call for every 20 seconds that
// passes.
//
// The fields are deliberately generic so that more than one algorithm (token
// bucket, leaky bucket, window-based designs) can be configured the same way.
// Your design does not have to use every field.
type Config struct {
	// Capacity is X: the most calls a key may make per Interval.
	Capacity int
	// Interval is A: the period Capacity applies to.
	Interval time.Duration
	// RefillAmount is Y: how many calls come back per RefillInterval.
	RefillAmount int
	// RefillInterval is B: how often RefillAmount calls come back.
	RefillInterval time.Duration
}

// ErrInvalidConfig is returned (wrapped) by New when the Config cannot work.
var ErrInvalidConfig = errors.New("ratelimiter: invalid config")

// Limiter tracks calls per key. It is shared by every HTTP request the server
// handles, so Allow is called from many goroutines at once.
type Limiter struct {
	cfg   Config
	clock clock.Clock

	// YOUR CODE HERE: whatever per-key state your design needs.
}

// New returns a Limiter for cfg. A nil clk means clock.SystemClock.
func New(cfg Config, clk clock.Clock) (*Limiter, error) {
	if clk == nil {
		clk = clock.SystemClock{}
	}

	// YOUR CODE HERE: validate cfg (return an error wrapping
	// ErrInvalidConfig) and initialise any state the Limiter needs.

	return &Limiter{cfg: cfg, clock: clk}, nil
}

// Allow reports whether key may make a call right now. When it returns true
// the call has been counted against key's limit.
func (l *Limiter) Allow(key string) bool {
	// YOUR CODE HERE
	panic("ratelimiter: Allow not implemented")
}
