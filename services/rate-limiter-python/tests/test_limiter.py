import pytest

from rate_limiter import Config, RateLimiter


class FakeClock:
    """Only moves when the test calls advance()."""

    def __init__(self) -> None:
        self._now = 1_000.0

    def now(self) -> float:
        return self._now

    def advance(self, seconds: float) -> None:
        self._now += seconds


# 3 calls per minute, refilling 1 call every 20 seconds.
CONFIG = Config(capacity=3, interval=60, refill_amount=1, refill_interval=20)

# Several users so that state leaking between keys shows up as a failure.
USERS = ["alice", "bob", "carol"]


@pytest.fixture
def clock() -> FakeClock:
    return FakeClock()


@pytest.fixture
def limiter(clock: FakeClock) -> RateLimiter:
    return RateLimiter(CONFIG, clock)


def test_limiter_drain_limit_and_refill_per_user(limiter: RateLimiter, clock: FakeClock) -> None:
    # 1. Each user can spend their full capacity.
    for user in USERS:
        for call in range(1, CONFIG.capacity + 1):
            assert limiter.allow(user), f"{user}: call {call} of {CONFIG.capacity} was denied, want allowed"

    # 2. Each user is denied once over the limit.
    for user in USERS:
        assert not limiter.allow(user), (
            f"{user}: call {CONFIG.capacity + 1} was allowed, want denied (capacity is {CONFIG.capacity})"
        )

    clock.advance(CONFIG.refill_interval)

    # 3. Each user gets one call back after one refill interval.
    for user in USERS:
        assert limiter.allow(user), f"{user}: first call after {CONFIG.refill_interval}s was denied, want allowed"

    # 4. Each user is denied again after spending the refill.
    for user in USERS:
        assert not limiter.allow(user), f"{user}: second call after {CONFIG.refill_interval}s was allowed, want denied"
