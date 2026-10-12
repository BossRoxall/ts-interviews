import time
from typing import Protocol


class Clock(Protocol):
    """The limiter's only source of time.

    Production code uses SystemClock; tests use a fake clock they can advance
    by hand.
    """

    def now(self) -> float:
        """Current time in seconds. Only differences between values matter."""
        ...


class SystemClock:
    """Reads a monotonic clock."""

    def now(self) -> float:
        return time.monotonic()
