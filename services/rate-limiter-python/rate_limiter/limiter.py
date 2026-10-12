"""Limits how often each key (for example a user ID) may make a call."""

from dataclasses import dataclass

from rate_limiter.clock import Clock, SystemClock


@dataclass(frozen=True)
class Config:
    """The limit applied to every key.

    Each key may make up to ``capacity`` calls per ``interval``, and spent calls
    come back at ``refill_amount`` calls every ``refill_interval``. For example::

        Config(capacity=3, interval=60, refill_amount=1, refill_interval=20)

    allows a burst of 3 calls, then one more call for every 20 seconds that
    passes.

    The fields are deliberately generic so that more than one algorithm (token
    bucket, leaky bucket, window-based designs) can be configured the same way.
    Your design does not have to use every field.
    """

    capacity: int
    """X: the most calls a key may make per ``interval``."""
    interval: float
    """A: the period ``capacity`` applies to, in seconds."""
    refill_amount: int
    """Y: how many calls come back per ``refill_interval``."""
    refill_interval: float
    """B: how often ``refill_amount`` calls come back, in seconds."""


class InvalidConfigError(ValueError):
    """Raised by RateLimiter when the Config cannot work."""


class RateLimiter:
    """Tracks calls per key.

    It is shared by every HTTP request the server handles, so ``allow`` is
    called from many threads at once.
    """

    def __init__(self, config: Config, clock: Clock | None = None) -> None:
        self._config = config
        self._clock = clock or SystemClock()

        # YOUR CODE HERE: validate config (raise InvalidConfigError) and
        # initialize whatever per-key state your design needs.

    def allow(self, key: str) -> bool:
        """Report whether ``key`` may make a call right now.

        When it returns True the call has been counted against ``key``'s limit.
        """
        # YOUR CODE HERE
        raise NotImplementedError("RateLimiter.allow is not implemented")
