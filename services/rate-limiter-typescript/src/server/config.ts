/**
 * Everything the server reads from the environment.
 *
 *     PORT                                8080
 *     RATE_LIMIT_CAPACITY                 100
 *     RATE_LIMIT_INTERVAL_SECONDS         60
 *     RATE_LIMIT_REFILL_AMOUNT            10
 *     RATE_LIMIT_REFILL_INTERVAL_SECONDS  6
 */

import type { Config } from "../limiter.ts";
import { log } from "./log.ts";

export interface ServerConfig {
  port: number;
  rateLimit: Config;
}

export function configFromEnv(): ServerConfig {
  return {
    port: envInt("PORT", 8080),
    rateLimit: {
      capacity: envInt("RATE_LIMIT_CAPACITY", 100),
      intervalMs: envInt("RATE_LIMIT_INTERVAL_SECONDS", 60) * 1000,
      refillAmount: envInt("RATE_LIMIT_REFILL_AMOUNT", 10),
      refillIntervalMs: envInt("RATE_LIMIT_REFILL_INTERVAL_SECONDS", 6) * 1000,
    },
  };
}

function envInt(name: string, fallback: number): number {
  const raw = process.env[name] ?? "";
  if (raw === "") return fallback;
  const value = Number(raw);
  if (!Number.isInteger(value)) {
    log("ERROR", `${name}=${JSON.stringify(raw)} is not an integer`);
    process.exit(1);
  }
  return value;
}
