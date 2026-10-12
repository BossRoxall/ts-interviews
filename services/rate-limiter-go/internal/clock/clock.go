// Package clock provides the time source used by the rate limiter.
package clock

import "time"

// Clock is the limiter's only source of time. Production code uses
// SystemClock; tests use a fake clock they can advance by hand.
type Clock interface {
	Now() time.Time
}

// SystemClock reads the wall clock.
type SystemClock struct{}

// Now returns the current time.
func (SystemClock) Now() time.Time { return time.Now() }
