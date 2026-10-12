"""Every endpoint the service exposes, keyed by (method, path)."""

from collections.abc import Callable

from rate_limiter.server import handlers

Route = Callable[[handlers.Request], handlers.Response]

ROUTES: dict[tuple[str, str], Route] = {
    ("GET", "/healthz"): handlers.healthz,  # liveness
    ("GET", "/readyz"): handlers.readyz,  # readiness
    ("POST", "/v1/check"): handlers.check,  # rate limit check
}
