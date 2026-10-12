"""Request handlers. Each takes a Request and returns a (status, JSON body) pair."""

import json
import logging
import threading
from dataclasses import dataclass
from http import HTTPStatus
from typing import Any

from rate_limiter.limiter import RateLimiter

log = logging.getLogger("rate_limiter")


@dataclass(frozen=True)
class Request:
    body: bytes
    limiter: RateLimiter
    ready: threading.Event


Response = tuple[HTTPStatus, dict[str, Any]]


def healthz(request: Request) -> Response:
    """Report that the process is up."""
    return HTTPStatus.OK, {"status": "ok"}


def readyz(request: Request) -> Response:
    """Report whether the server is accepting traffic.

    Returns 503 while shutting down so the pod is taken out of the Service first.
    """
    if request.ready.is_set():
        return HTTPStatus.OK, {"status": "ready"}
    return HTTPStatus.SERVICE_UNAVAILABLE, {"status": "not ready"}


def check(request: Request) -> Response:
    """Answer {"key": "alice"} with 200 {"allowed": true} or 429 {"allowed": false}."""
    try:
        key = json.loads(request.body or b"{}").get("key")
    except (ValueError, AttributeError):
        key = None
    if not isinstance(key, str) or key == "":
        return HTTPStatus.BAD_REQUEST, {"error": 'body must be JSON like {"key": "alice"}'}

    try:
        allowed = request.limiter.allow(key)
    except Exception as exc:  # e.g. the unimplemented limiter
        log.error("handler error on /v1/check: %r", exc)
        return HTTPStatus.INTERNAL_SERVER_ERROR, {"error": str(exc)}

    if allowed:
        return HTTPStatus.OK, {"allowed": True}
    return HTTPStatus.TOO_MANY_REQUESTS, {"allowed": False}
