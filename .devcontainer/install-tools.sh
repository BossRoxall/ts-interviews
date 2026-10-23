#!/usr/bin/env bash
# Installs the pinned toolchain for this repo into /usr/local on Debian or
# Ubuntu (amd64 or arm64). The dev container image runs it at build time, and
# SANDBOX.md uses it for a throwaway VM.
#
# It is meant for disposable environments. On your own laptop, use the package
# manager commands in README.md instead.
#
# Every download is checked against the SHA-256 pinned below before anything
# is unpacked.
set -euo pipefail

GO_VERSION=1.27.3
GOPLS_VERSION=0.24.0
NODE_VERSION=26.13.0
UV_VERSION=0.13.7
PYTHON_VERSION=3.13
KUBECTL_VERSION=1.37.2
HELM_VERSION=4.3.1
KUSTOMIZE_VERSION=5.8.3
KIND_VERSION=0.33.0
TILT_VERSION=0.37.8

case "$(uname -m)" in
  x86_64 | amd64)
    ARCH=amd64
    GO_SHA=65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40
    NODE_ARCH=x64
    NODE_SHA=d28a96af6c38596f6d9d97c0ee531c13d238d757742387d1928de822d959e997
    UV_ARCH=x86_64
    UV_SHA=c5221ff41b7892bfc223543883600125f6a47fe845093e0f5bdff4f359196d83
    KUBECTL_SHA=c04487523e7aa7d1493f7808174c899db93d7a9cf57531420cd76b9bd415f8f4
    HELM_SHA=0403a88ae20913d11b10a2467b50087721fea7581736da6b9d3992f5357fe54b
    KUSTOMIZE_SHA=cb9e31198d3f63b44848bd0afc4c8efc56e6cfef9c93a4a39a211940e6decd61
    KIND_SHA=aee6151561422756b764a4ae28e7f44cda5af5a9eead3cc9985112b1de8d8e0d
    TILT_ARCH=x86_64
    TILT_SHA=e39513f958705d30524d55f68b1d46fbef2b46a0da568efd00e6383dee6be44f
    ;;
  aarch64 | arm64)
    ARCH=arm64
    GO_SHA=a06e99309dbcbdca1a8ccd2282c3ee146e55930886e26b4f1193269587165f49
    NODE_ARCH=arm64
    NODE_SHA=406fe3fbde1c054adff8b297c16fdb5af6659db327c1ac8acd46b56ef1af9a92
    UV_ARCH=aarch64
    UV_SHA=71279f8d3c67241e807b9d8db872d1bf541579fb9f08b7d1eb1682196d1e57dc
    KUBECTL_SHA=fcc16f63186659d9722014b7e02d5fdd2d4e669cb8c0ab66365d7d3b8ba03509
    HELM_SHA=8cf8f66ecbc252ac8f7f06cca68aba88508ddfde3a6a34feb1678ad1694b4f88
    KUSTOMIZE_SHA=9867b76482cfc6d25b546abc1d2b5a32bd137505aba640649b4b378f64f6d68f
    KIND_SHA=20022bee6cfcd5086cb7234d218e3454e6090022f2a8f55d1fa7fcf42c3867a2
    TILT_ARCH=arm64
    TILT_SHA=7a21281e830930be158202b7f8405ba62445613788b4d9708512ffe51f5cce90
    ;;
  *)
    echo "unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

SUDO=""
if [[ "$(id -u)" -ne 0 ]]; then
  SUDO="sudo"
fi

BIN=/usr/local/bin
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

# fetch URL SHA256 FILENAME: download into $WORK and verify the checksum.
fetch() {
  local url="$1" sha="$2" name="$3"
  echo "fetching ${name}"
  curl -fsSL --retry 3 -o "$WORK/$name" "$url"
  echo "${sha}  $WORK/$name" | sha256sum -c --quiet -
}

echo "installing base packages"
$SUDO apt-get update -qq
$SUDO env DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends \
  ca-certificates curl git make tar xz-utils >/dev/null

fetch "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" "$GO_SHA" go.tar.gz
$SUDO rm -rf /usr/local/go
$SUDO tar -C /usr/local -xzf "$WORK/go.tar.gz"
$SUDO ln -sf /usr/local/go/bin/go /usr/local/go/bin/gofmt "$BIN/"

# gopls is the language server the VS Code Go extension uses. `go install`
# checks the module against the Go checksum database (sum.golang.org).
echo "installing gopls ${GOPLS_VERSION}"
$SUDO env GOBIN="$BIN" GOPATH="$WORK/gopath" GOCACHE="$WORK/gocache" GOFLAGS=-modcacherw \
  /usr/local/go/bin/go install "golang.org/x/tools/gopls@v${GOPLS_VERSION}"

fetch "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-${NODE_ARCH}.tar.xz" "$NODE_SHA" node.tar.xz
$SUDO rm -rf /usr/local/lib/nodejs
$SUDO mkdir -p /usr/local/lib/nodejs
$SUDO tar -C /usr/local/lib/nodejs --strip-components=1 -xJf "$WORK/node.tar.xz"
for tool in node npm npx; do
  $SUDO ln -sf "/usr/local/lib/nodejs/bin/$tool" "$BIN/$tool"
done

fetch "https://github.com/astral-sh/uv/releases/download/${UV_VERSION}/uv-${UV_ARCH}-unknown-linux-gnu.tar.gz" "$UV_SHA" uv.tar.gz
tar -C "$WORK" -xzf "$WORK/uv.tar.gz"
$SUDO install -m 0755 "$WORK/uv-${UV_ARCH}-unknown-linux-gnu/uv" "$WORK/uv-${UV_ARCH}-unknown-linux-gnu/uvx" "$BIN/"

# uv pins the Python builds it knows about and verifies their hashes itself.
echo "installing python ${PYTHON_VERSION} with uv"
$SUDO env UV_PYTHON_INSTALL_DIR=/usr/local/share/uv/python UV_PYTHON_BIN_DIR="$BIN" \
  "$BIN/uv" python install --quiet --default "$PYTHON_VERSION"

fetch "https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/${ARCH}/kubectl" "$KUBECTL_SHA" kubectl
$SUDO install -m 0755 "$WORK/kubectl" "$BIN/kubectl"

fetch "https://get.helm.sh/helm-v${HELM_VERSION}-linux-${ARCH}.tar.gz" "$HELM_SHA" helm.tar.gz
tar -C "$WORK" -xzf "$WORK/helm.tar.gz"
$SUDO install -m 0755 "$WORK/linux-${ARCH}/helm" "$BIN/helm"

fetch "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_linux_${ARCH}.tar.gz" "$KUSTOMIZE_SHA" kustomize.tar.gz
tar -C "$WORK" -xzf "$WORK/kustomize.tar.gz"
$SUDO install -m 0755 "$WORK/kustomize" "$BIN/kustomize"

fetch "https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-linux-${ARCH}" "$KIND_SHA" kind
$SUDO install -m 0755 "$WORK/kind" "$BIN/kind"

fetch "https://github.com/tilt-dev/tilt/releases/download/v${TILT_VERSION}/tilt.${TILT_VERSION}.linux.${TILT_ARCH}.tar.gz" "$TILT_SHA" tilt.tar.gz
tar -C "$WORK" -xzf "$WORK/tilt.tar.gz" tilt
$SUDO install -m 0755 "$WORK/tilt" "$BIN/tilt"

echo
echo "installed:"
go version
gopls version | awk "NR == 1"
echo "node $(node --version)"
uv --version
python3 --version
kubectl version --client | awk "NR == 1"
helm version --short
echo "kustomize $(kustomize version)"
kind version
tilt version
