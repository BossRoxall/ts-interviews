/** Limits how often each key (for example a user ID) may make a call. */

import { type Clock, SystemClock } from "./clock.ts";

/**
 * The limit applied to every key.
 *
 * Each key may make up to `capacity` calls per `intervalMs`, and spent calls
 * come back at `refillAmount` calls every `refillIntervalMs`. For example:
 *
 *     { capacity: 3, intervalMs: 60_000, refillAmount: 1, refillIntervalMs: 20_000 }
 *
 * allows a burst of 3 calls, then one more call for every 20 seconds that
 * passes.
 *
 * The fields are deliberately generic so that more than one algorithm (token
 * bucket, leaky bucket, window-based designs) can be configured the same way.
 * Your design does not have to use every field.
 */
export interface Config {
  /** X: the most calls a key may make per `intervalMs`. */
  capacity: number;
  /** A: the period `capacity` applies to, in milliseconds. */
  intervalMs: number;
  /** Y: how many calls come back per `refillIntervalMs`. */
  refillAmount: number;
  /** B: how often `refillAmount` calls come back, in milliseconds. */
  refillIntervalMs: number;
}

/** Thrown by RateLimiter when the Config cannot work. */
export class InvalidConfigError extends Error {
  override name = "InvalidConfigError";
}

/**
 * Tracks calls per key. It is shared by every HTTP request the server
 * handles, so `allow` is called for many concurrent requests.
 */
export class RateLimiter {
  private readonly config: Config;
  private readonly clock: Clock;

  // YOUR CODE HERE: whatever per-key state your design needs.

  constructor(config: Config, clock: Clock = new SystemClock()) {
    this.config = config;
    this.clock = clock;

    // YOUR CODE HERE: validate config (throw InvalidConfigError) and
    // initialize any state the limiter needs.
  }

  /**
   * Report whether `key` may make a call right now. When it returns true the
   * call has been counted against `key`'s limit.
   */
  allow(key: string): boolean {
    // YOUR CODE HERE
    throw new Error("RateLimiter.allow is not implemented");
  }
}
