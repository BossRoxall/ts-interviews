# Rate limiter (Go)

The same exercise exists in Go, Python and TypeScript. Pick whichever language you are most comfortable in; they all have the same API, tests and HTTP server.

## The task

Implement a per-key rate limiter. A key is usually a user ID, and each key gets its own independent limit.

The limiter is configured with:

| Field                | Meaning                                       |
| -------------------- | --------------------------------------------- |
| `Capacity` (X)       | the most calls a key may make per `Interval`  |
| `Interval` (A)       | the period `Capacity` applies to              |
| `RefillAmount` (Y)   | how many calls come back per `RefillInterval` |
| `RefillInterval` (B) | how often `RefillAmount` calls come back      |

So `Capacity: 3, Interval: 1m, RefillAmount: 1, RefillInterval: 20s` means a user can make 3 calls straight away, is then refused, and gets one more call for every 20 seconds that passes.

The config is generic on purpose. It should be able to describe a token bucket, a leaky bucket or a window-based design, so you choose the algorithm. Your design does not have to use every field; tell us which ones it uses and why.

Fill in the `YOUR CODE HERE` sections in `limiter.go`:

1. the per-key state your design needs, on `Limiter`
2. config validation and setup in `New`
3. `Allow(key)`, which returns `true` and counts the call if the key is within its limit, and `false` otherwise

Things to keep in mind:

- Time comes from the injected `clock.Clock` (in `internal/clock`). Do not call `time.Now()` directly, or the tests cannot control time.
- The HTTP server in `internal/server` shares one `Limiter` across all requests, and Go's HTTP server handles requests on many goroutines at once.
- You don't need to change the tests or the server.

Talk us through your thinking as you go. We care more about how you reason about the trade-offs than about finishing every stretch goal.

## Running the tests

```sh
go test ./...
```

Run the tests by package as above (or `go test .`), not by file: `go test ./limiter_test.go` compiles that one file on its own and fails with `undefined: Config`.

The tests use a fake clock and three users (`alice`, `bob`, `carol`). They check that each user can spend their full capacity, that each user is refused once over the limit, that each user gets exactly one call back after one refill interval, and that they are refused again after spending it. The tests fail with `not implemented` until `Allow` is written.

## Stretch goals

If you finish early, pick any of these to discuss or build:

- Return a decision with the remaining calls and how long until the next call is allowed, so the server can send a `Retry-After` header.
- Memory: a new key creates new state forever. How would you stop idle keys from growing memory without bound?
- Add a second algorithm behind the same `Config`.
- The service runs several replicas behind a load balancer. What happens to each user's limit, and what would you change?

## Running the server

```sh
go run ./cmd/server

curl -s localhost:8080/healthz
curl -s localhost:8080/readyz
curl -s -XPOST localhost:8080/v1/check -d '{"key": "alice"}'
```

Endpoints are listed in `internal/server/routes.go`. `POST /v1/check` answers `200 {"allowed": true}` or `429 {"allowed": false}`. `GET /healthz` reports that the process is up, and `GET /readyz` reports that it is accepting traffic (it returns 503 while shutting down).

Configuration comes from environment variables:

| Variable                             | Default |
| ------------------------------------ | ------- |
| `PORT`                               | `8080`  |
| `RATE_LIMIT_CAPACITY`                | `100`   |
| `RATE_LIMIT_INTERVAL_SECONDS`        | `60`    |
| `RATE_LIMIT_REFILL_AMOUNT`           | `10`    |
| `RATE_LIMIT_REFILL_INTERVAL_SECONDS` | `6`     |

This is the version production runs; `Dockerfile` builds its image. `tilt up` from the repo root deploys it to a local cluster (see the root README).
