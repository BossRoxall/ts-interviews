/** Request handlers. Each takes a Request and returns a status and JSON body. */

import type { RateLimiter } from "../limiter.ts";
import { log } from "./log.ts";

export interface Request {
  body: string;
  limiter: RateLimiter;
  isReady: () => boolean;
}

export interface Response {
  status: number;
  body: unknown;
}

export type Handler = (request: Request) => Response;

/** Report that the process is up. */
export const healthz: Handler = () => ({ status: 200, body: { status: "ok" } });

/**
 * Report whether the server is accepting traffic. Returns 503 while shutting
 * down so the pod is taken out of the Service first.
 */
export const readyz: Handler = ({ isReady }) =>
  isReady() ? { status: 200, body: { status: "ready" } } : { status: 503, body: { status: "not ready" } };

/** Answer {"key": "alice"} with 200 {"allowed": true} or 429 {"allowed": false}. */
export const check: Handler = ({ body, limiter }) => {
  let key: unknown;
  try {
    key = JSON.parse(body || "{}").key;
  } catch {
    key = undefined;
  }
  if (typeof key !== "string" || key === "") {
    return { status: 400, body: { error: 'body must be JSON like {"key": "alice"}' } };
  }

  try {
    const allowed = limiter.allow(key);
    return { status: allowed ? 200 : 429, body: { allowed } };
  } catch (err) {
    // e.g. the unimplemented limiter
    const message = err instanceof Error ? err.message : String(err);
    log("ERROR", "handler error on /v1/check", { error: message });
    return { status: 500, body: { error: message } };
  }
};
