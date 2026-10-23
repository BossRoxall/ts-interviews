#!/usr/bin/env bash
# Reports which tools are installed and which exercises they cover.
# Read-only: it installs nothing and changes nothing.
set -uo pipefail

ok=0
missing=0

check() {
  local tool="$1" purpose="$2" version_cmd="$3"
  if command -v "$tool" >/dev/null 2>&1; then
    printf '  %-10s %-48s %s\n' "$tool" "$(eval "$version_cmd" 2>/dev/null | awk 'NR == 1')" "$purpose"
    ok=$((ok + 1))
  else
    printf '  %-10s %-48s %s\n' "$tool" "MISSING" "$purpose"
    missing=$((missing + 1))
  fi
}

echo "Language toolchains (you only need the one you will use):"
check go      "Go service, loadgen"          "go version"
check uv      "Python service"               "uv --version"
check node    "TypeScript service"           "node --version"
echo
echo "Running the service in a local cluster:"
check docker    "building images"            "docker --version"
check kubectl   "talking to the cluster"     "kubectl version --client"
check helm      "rendering the chart"        "helm version --short"
check kustomize "rendering the overlay"      "kustomize version"
check tilt      "tilt up"                    "tilt version"
echo

if command -v kubectl >/dev/null 2>&1; then
  context="$(kubectl config current-context 2>/dev/null || true)"
  if [[ -z "$context" ]]; then
    echo "No current kube context. Start a local Kubernetes engine (see README.md)."
  elif kubectl version --request-timeout=3s >/dev/null 2>&1; then
    echo "Kube context '$context' is reachable."
  else
    echo "Kube context '$context' is set but not reachable. Is your local engine running?"
  fi
fi

echo
echo "$ok found, $missing missing."
