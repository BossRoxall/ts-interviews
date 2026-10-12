/** Every endpoint the service exposes, keyed by "METHOD path". */

import { check, type Handler, healthz, readyz } from "./handlers.ts";

export const routes: Record<string, Handler> = {
  "GET /healthz": healthz, // liveness
  "GET /readyz": readyz, // readiness
  "POST /v1/check": check, // rate limit check
};
