# Rate limiter (Python)

The same exercise exists in Go, Python and TypeScript. Pick whichever language you are most comfortable in; they all have the same API, tests and HTTP server.

## The task

Implement a per-key rate limiter. A key is usually a user ID, and each key gets its own independent limit.

The limiter is configured with:

| Field                 | Meaning                                               |
| --------------------- | ----------------------------------------------------- |
| `capacity` (X)        | the most calls a key may make per `interval`          |
| `interval` (A)        | the period `capacity` applies to, in seconds          |
| `refill_amount` (Y)   | how many calls come back per `refill_interval`        |
| `refill_interval` (B) | how often `refill_amount` calls come back, in seconds |

So `capacity=3, interval=60, refill_amount=1, refill_interval=20` means a user can make 3 calls straight away, is then refused, and gets one more call for every 20 seconds that passes.

The config is generic on purpose. It should be able to describe a token bucket, a leaky bucket or a window-based design, so you choose the algorithm. Your design does not have to use every field; tell us which ones it uses and why.

Fill in the `YOUR CODE HERE` sections in `rate_limiter/limiter.py`:

1. the per-key state your design needs, on `RateLimiter`
2. config validation and setup in `__init__`
3. `allow(key)`, which returns `True` and counts the call if the key is within its limit, and `False` otherwise

Things to keep in mind:

- Time comes from the injected `Clock`. Do not call `time.time()` or `time.monotonic()` directly, or the tests cannot control time.
- The HTTP server in `rate_limiter/server/` shares one `RateLimiter` across all requests, and `ThreadingHTTPServer` handles each request on its own thread.
- You don't need to change the tests or the server.

Talk us through your thinking as you go. We care more about how you reason about the trade-offs than about finishing every stretch goal.

## Running the tests

```sh
uv run pytest
```

Without uv: `python -m venv .venv && . .venv/bin/activate && pip install pytest && pytest`.

The tests use a fake clock and three users (`alice`, `bob`, `carol`). They check that each user can spend their full capacity, that each user is refused once over the limit, that each user gets exactly one call back after one refill interval, and that they are refused again after spending it. The tests fail with `NotImplementedError` until `allow` is written.

## Stretch goals

If you finish early, pick any of these to discuss or build:

- Return a decision with the remaining calls and how long until the next call is allowed, so the server can send a `Retry-After` header.
- Memory: a new key creates new state forever. How would you stop idle keys from growing memory without bound?
- Add a second algorithm behind the same `Config`.
- The service runs several replicas behind a load balancer. What happens to each user's limit, and what would you change?

## Running the server

```sh
uv run python -m rate_limiter.server

curl -s localhost:8080/healthz
curl -s localhost:8080/readyz
curl -s -XPOST localhost:8080/v1/check -d '{"key": "alice"}'
```

Endpoints are listed in `rate_limiter/server/routes.py`. `POST /v1/check` answers `200 {"allowed": true}` or `429 {"allowed": false}`. `GET /healthz` reports that the process is up, and `GET /readyz` reports that it is accepting traffic (it returns 503 while shutting down).

Configuration comes from environment variables:

| Variable                             | Default |
| ------------------------------------ | ------- |
| `PORT`                               | `8080`  |
| `RATE_LIMIT_CAPACITY`                | `100`   |
| `RATE_LIMIT_INTERVAL_SECONDS`        | `60`    |
| `RATE_LIMIT_REFILL_AMOUNT`           | `10`    |
| `RATE_LIMIT_REFILL_INTERVAL_SECONDS` | `6`     |

Production runs the Go version. To run this one in a local cluster instead, set `"selected_language": "python"` in the root `tilt_config.json` and run `tilt up` (see the root README).
