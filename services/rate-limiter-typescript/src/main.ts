/**
 * Runs the rate limiter HTTP service. Routes are listed in
 * src/server/routes.ts and configuration in src/server/config.ts.
 */

import { InvalidConfigError } from "./limiter.ts";
import { configFromEnv } from "./server/config.ts";
import { log } from "./server/log.ts";
import { Server } from "./server/server.ts";

let server: Server;
try {
  server = new Server(configFromEnv());
} catch (err) {
  if (err instanceof InvalidConfigError) {
    log("ERROR", "invalid rate limit config", { error: err.message });
    process.exit(1);
  }
  throw err;
}

process.on("SIGTERM", () => server.shutdown());
process.on("SIGINT", () => server.shutdown());

server.listen();
