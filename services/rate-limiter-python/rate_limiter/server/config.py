"""Everything the server reads from the environment.

    PORT                                8080
    RATE_LIMIT_CAPACITY                 100
    RATE_LIMIT_INTERVAL_SECONDS         60
    RATE_LIMIT_REFILL_AMOUNT            10
    RATE_LIMIT_REFILL_INTERVAL_SECONDS  6
"""

import os
from dataclasses import dataclass

from rate_limiter.limiter import Config


@dataclass(frozen=True)
class ServerConfig:
    port: int
    rate_limit: Config

    @classmethod
    def from_env(cls) -> "ServerConfig":
        return cls(
            port=_env_int("PORT", 8080),
            rate_limit=Config(
                capacity=_env_int("RATE_LIMIT_CAPACITY", 100),
                interval=_env_int("RATE_LIMIT_INTERVAL_SECONDS", 60),
                refill_amount=_env_int("RATE_LIMIT_REFILL_AMOUNT", 10),
                refill_interval=_env_int("RATE_LIMIT_REFILL_INTERVAL_SECONDS", 6),
            ),
        )


def _env_int(name: str, fallback: int) -> int:
    raw = os.environ.get(name, "")
    if raw == "":
        return fallback
    try:
        return int(raw)
    except ValueError:
        raise SystemExit(f"{name}={raw!r} is not an integer") from None
