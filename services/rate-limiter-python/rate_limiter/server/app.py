"""HTTP server lifecycle: listening, dispatching to routes, graceful shutdown."""

import json
import logging
import signal
import threading
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Any

from rate_limiter.limiter import RateLimiter
from rate_limiter.server.config import ServerConfig
from rate_limiter.server.handlers import Request
from rate_limiter.server.routes import ROUTES

log = logging.getLogger("rate_limiter")


class Server:
    def __init__(self, config: ServerConfig) -> None:
        self.config = config
        self.limiter = RateLimiter(config.rate_limit)
        self.ready = threading.Event()

    def run(self) -> None:
        httpd = ThreadingHTTPServer(("", self.config.port), self._handler_class())

        def shutdown(signum: int, _frame: object) -> None:
            log.info("shutting down")
            self.ready.clear()
            threading.Thread(target=httpd.shutdown, daemon=True).start()

        signal.signal(signal.SIGTERM, shutdown)
        signal.signal(signal.SIGINT, shutdown)

        log.info("listening on port %d with %s", self.config.port, self.config.rate_limit)
        self.ready.set()
        httpd.serve_forever()

    def _handler_class(self) -> type[BaseHTTPRequestHandler]:
        server = self

        class Handler(BaseHTTPRequestHandler):
            def do_GET(self) -> None:
                self._dispatch("GET")

            def do_POST(self) -> None:
                self._dispatch("POST")

            def _dispatch(self, method: str) -> None:
                route = ROUTES.get((method, self.path))
                if route is None:
                    self._send(HTTPStatus.NOT_FOUND, {"error": "not found"})
                    return
                length = min(int(self.headers.get("Content-Length", "0") or 0), 4096)
                request = Request(body=self.rfile.read(length), limiter=server.limiter, ready=server.ready)
                self._send(*route(request))

            def _send(self, status: HTTPStatus, body: dict[str, Any]) -> None:
                payload = json.dumps(body).encode()
                self.send_response(status)
                self.send_header("Content-Type", "application/json")
                self.send_header("Content-Length", str(len(payload)))
                self.end_headers()
                self.wfile.write(payload)

            def log_message(self, format: str, *args: Any) -> None:  # quieter default logging
                pass

        return Handler
