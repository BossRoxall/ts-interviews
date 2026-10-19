# rate-limiter

A small per-key rate limiting service and the infrastructure that deploys it.

```
services/
  rate-limiter-go/          Go implementation (the one production runs)
  rate-limiter-python/      Python implementation
  rate-limiter-typescript/  TypeScript implementation
infra/
  charts/rate-limiter/      Helm chart
  deployments/production/   production overlay (Kustomize + Helm values)
  argo/                     Argo CD Application for production
tools/loadgen/              load generator for local testing
Tiltfile                    local deployment with Tilt
```

All three services expose the same HTTP API and share the same tests. Each service has its own README with setup and test instructions.

## Prerequisites

| Tool | Used for |
| --- | --- |
| Go 1.27+ | Go service, load generator |
| Python 3.12+ and [uv](https://docs.astral.sh/uv/) | Python service |
| Node 23.6+ | TypeScript service |
| Helm 3.8+ and Kustomize 5+ | rendering the chart and overlay |
| Tilt and a local Kubernetes engine | running the production overlay locally |

Any local engine works: OrbStack, Docker Desktop, minikube, kind, Rancher Desktop or Colima. The Tiltfile refuses to run against any other kube context.

## Common tasks

```sh
make test     # run the tests for all three services
make lint     # go vet, tsc, helm lint
make render   # print the production manifests
make up       # tilt up
make down     # tilt down
```

## Running locally with Tilt

`tilt up` builds one of the services, renders `infra/deployments/production` with `kustomize build --enable-helm`, applies it to your current local cluster and port-forwards the pod to `localhost:8080`. Press space to open the Tilt UI.

```sh
curl -s -XPOST localhost:8080/v1/check -d '{"key": "alice"}'
```

`tilt_config.json` holds the Tiltfile options:

| Key | Default | Meaning |
| --- | --- | --- |
| `selected_language` | `go` | which service to build and deploy: `go`, `python`, or `typescript` (`node` and `nodejs` also work) |
| `host_port` | `8080` | local port forwarded to the pod |
| `loadgen` | `true` | show the manual `loadgen` resource |

Whichever language you pick is built under the production image name, so the manifests are identical either way. Override an option for one run with `tilt up -- --selected_language python`, or change it on a running Tilt with `tilt args -- --selected_language python`.

The `loadgen` resource in the Tilt UI is manual: click it to send a burst of requests, each from a different key, through the port-forward.

## Deployment

CI builds and pushes `registry.gitlab.com/turionspace/interviews/rate-limiter` from `services/rate-limiter-go`. Argo CD watches `infra/deployments/production` on `main` and syncs it to the production cluster.
