import assert from "node:assert/strict";
import { test } from "node:test";

import type { Clock } from "../src/clock.ts";
import { type Config, RateLimiter } from "../src/limiter.ts";

/** Only moves when the test calls advance(). */
class FakeClock implements Clock {
  private current = 1_000_000;

  now(): number {
    return this.current;
  }

  advance(ms: number): void {
    this.current += ms;
  }
}

// 3 calls per minute, refilling 1 call every 20 seconds.
const config: Config = {
  capacity: 3,
  intervalMs: 60_000,
  refillAmount: 1,
  refillIntervalMs: 20_000,
};

// Several users so that state leaking between keys shows up as a failure.
const users = ["alice", "bob", "carol"];

test("limiter drains, limits and refills each user independently", async (t) => {
  const clock = new FakeClock();
  const limiter = new RateLimiter(config, clock);

  await t.test("1 each user can spend their full capacity", () => {
    for (const user of users) {
      for (let call = 1; call <= config.capacity; call++) {
        assert.ok(limiter.allow(user), `${user}: call ${call} of ${config.capacity} was denied, want allowed`);
      }
    }
  });

  await t.test("2 each user is denied once over the limit", () => {
    for (const user of users) {
      assert.ok(
        !limiter.allow(user),
        `${user}: call ${config.capacity + 1} was allowed, want denied (capacity is ${config.capacity})`,
      );
    }
  });

  clock.advance(config.refillIntervalMs);

  await t.test("3 each user gets one call back after one refill interval", () => {
    for (const user of users) {
      assert.ok(limiter.allow(user), `${user}: first call after ${config.refillIntervalMs}ms was denied, want allowed`);
    }
  });

  await t.test("4 each user is denied again after spending the refill", () => {
    for (const user of users) {
      assert.ok(!limiter.allow(user), `${user}: second call after ${config.refillIntervalMs}ms was allowed, want denied`);
    }
  });
});
