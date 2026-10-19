.PHONY: test test-go test-python test-typescript lint render up down

GO_DIR := services/rate-limiter-go
PY_DIR := services/rate-limiter-python
TS_DIR := services/rate-limiter-typescript

test: test-go test-python test-typescript

test-go:
	cd $(GO_DIR) && go test ./...

test-python:
	cd $(PY_DIR) && uv run pytest

test-typescript:
	cd $(TS_DIR) && npm test

lint:
	cd $(GO_DIR) && go vet ./...
	cd tools/loadgen && go vet ./...
	cd $(TS_DIR) && npm install --silent && npm run typecheck
	helm lint infra/charts/rate-limiter

# Print the production manifests exactly as Argo CD and Tilt render them.
render:
	kustomize build --enable-helm infra/deployments/production

up:
	tilt up

down:
	tilt down
