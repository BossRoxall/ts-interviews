/** HTTP server lifecycle: listening, dispatching to routes, graceful shutdown. */

import { createServer, type IncomingMessage, type Server as HttpServer, type ServerResponse } from "node:http";

import { RateLimiter } from "../limiter.ts";
import type { ServerConfig } from "./config.ts";
import { log } from "./log.ts";
import { routes } from "./routes.ts";

export class Server {
  private readonly config: ServerConfig;
  private readonly limiter: RateLimiter;
  private readonly http: HttpServer;
  private ready = false;

  constructor(config: ServerConfig) {
    this.config = config;
    this.limiter = new RateLimiter(config.rateLimit);
    this.http = createServer((req, res) => void this.dispatch(req, res));
  }

  listen(): void {
    this.http.listen(this.config.port, () => {
      log("INFO", "listening", { port: this.config.port, rateLimit: this.config.rateLimit });
      this.ready = true;
    });
  }

  shutdown(): void {
    log("INFO", "shutting down");
    this.ready = false;
    this.http.close(() => process.exit(0));
  }

  private async dispatch(req: IncomingMessage, res: ServerResponse): Promise<void> {
    const route = routes[`${req.method} ${req.url}`];
    if (route === undefined) {
      send(res, 404, { error: "not found" });
      return;
    }
    const { status, body } = route({
      body: await readBody(req),
      limiter: this.limiter,
      isReady: () => this.ready,
    });
    send(res, status, body);
  }
}

function send(res: ServerResponse, status: number, body: unknown): void {
  const payload = JSON.stringify(body);
  res.writeHead(status, { "Content-Type": "application/json", "Content-Length": Buffer.byteLength(payload) });
  res.end(payload);
}

async function readBody(req: IncomingMessage, limit = 4096): Promise<string> {
  let body = "";
  for await (const chunk of req) {
    body += chunk;
    if (body.length > limit) break;
  }
  return body;
}
