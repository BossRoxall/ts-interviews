#!/usr/bin/env bash
# Starts (or creates) a one-node kind cluster inside the dev container so
# `tilt up` has somewhere to deploy. Runs on every container start.
set -euo pipefail

CLUSTER=interview

for _ in $(seq 1 30); do
  docker info >/dev/null 2>&1 && break
  sleep 2
done
if ! docker info >/dev/null 2>&1; then
  echo "docker is not running inside the dev container; see SANDBOX.md" >&2
  exit 1
fi

if kind get clusters 2>/dev/null | grep -qx "$CLUSTER"; then
  docker start "${CLUSTER}-control-plane" >/dev/null 2>&1 || true
else
  kind create cluster --name "$CLUSTER" --wait 120s
fi

kind export kubeconfig --name "$CLUSTER" >/dev/null
kubectl config use-context "kind-${CLUSTER}" >/dev/null
echo "kind cluster '${CLUSTER}' is ready (context kind-${CLUSTER}). Run: tilt up"
