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
scripts/doctor.sh           checks which tools are installed
.devcontainer/              dev container with every tool preinstalled
Tiltfile                    local deployment with Tilt
```

All three services expose the same HTTP API and share the same tests. Each service has its own README with setup and test instructions.

## Getting set up

You only need what the task in front of you needs:

| To | You need |
| --- | --- |
| implement the limiter and run its tests | git and one language toolchain: Go 1.27+, Python 3.12+ with [uv](https://docs.astral.sh/uv/), or Node 23.6+ |
| run the service in a local cluster | Docker, a local Kubernetes engine, kubectl, Helm, Kustomize and [Tilt](https://tilt.dev) |
| review a merge request | a browser |

The quickest way to start, with nothing installed or run on your own machine, is a GitHub Codespace. Everything above comes preinstalled, with a local cluster ready for `tilt up`:

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new/BossRoxall/ts-interviews?quickstart=1)

[SANDBOX.md](SANDBOX.md) has the details, plus a dev container and a throwaway VM as alternatives.

### macOS

Install [Homebrew](https://brew.sh) if you don't have it, then:

```sh
xcode-select --install       # git and make, if macOS prompts for them

# The language you will use (one is enough)
brew install go              # Go
brew install uv              # Python
brew install node            # TypeScript

# Running in a local cluster
brew install kubectl helm kustomize tilt
brew install --cask orbstack # or Docker Desktop
```

Then turn on Kubernetes: in OrbStack it is under Settings > Kubernetes, and in Docker Desktop under Settings > Kubernetes > Enable Kubernetes. Both have license terms for commercial use; kind and minikube are free alternatives if that matters on your machine.

### Linux

Install Docker Engine and your language toolchain from your distribution's packages. For the cluster tooling, [kind](https://kind.sigs.k8s.io) gives you a local cluster with `kind create cluster`. On Debian or Ubuntu, `.devcontainer/install-tools.sh` installs every tool at the pinned versions with checksum verification, but it writes to `/usr/local`, so it is best run in a VM or container (see [SANDBOX.md](SANDBOX.md)).

### Windows

Use WSL 2 with Ubuntu and follow the Linux steps inside it, with Docker Desktop's WSL integration turned on. The dev container in [SANDBOX.md](SANDBOX.md) also works from VS Code on Windows.

### Checking your setup

```sh
make doctor     # or: bash scripts/doctor.sh
```

It lists each tool as found or missing, says which task needs it, and checks that your current kube context is reachable. It installs nothing.

The Tiltfile only deploys to local engines (OrbStack, Docker Desktop, minikube, kind, Rancher Desktop, Colima) and refuses any other kube context.

## Common tasks

```sh
make doctor   # check which tools are installed
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
