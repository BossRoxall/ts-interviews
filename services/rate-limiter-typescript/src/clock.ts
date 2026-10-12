/**
 * The limiter's only source of time. Production code uses SystemClock; tests
 * use a fake clock they can advance by hand.
 */
export interface Clock {
  /** Current time in milliseconds. Only differences between values matter. */
  now(): number;
}

/** Reads a monotonic clock. */
export class SystemClock implements Clock {
  now(): number {
    return performance.now();
  }
}
