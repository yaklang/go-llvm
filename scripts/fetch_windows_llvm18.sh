#!/usr/bin/env bash
# Download the prebuilt LLVM 18.1.3 + LLD static archives for Windows amd64 and
# unpack them into llvm/18.1.3/windows-amd64.
#
# The tree is not committed: it is 188 MB unpacked, which pushes the Go module
# over the 500 MiB source-tree limit enforced by module proxies. Keeping it in a
# release asset keeps `go get github.com/yaklang/go-llvm` working everywhere.
#
# Rebuild the tree from source with scripts/build_windows_llvm18.sh when the
# LLVM version or the build options change.
set -euo pipefail

VERSION="18.1.3"
TAG="llvm-${VERSION}-windows-amd64"
ASSET="windows-amd64-llvm18.tar.zst"
SHA256="bd5b39a330c20423397a53c4ddfdeacad6453d7bb354f2b0868697609f81059e"
REPO="yaklang/go-llvm"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST="${ROOT}/llvm/${VERSION}"
URL="${GO_LLVM_WINDOWS_URL:-https://github.com/${REPO}/releases/download/${TAG}/${ASSET}}"

if [ -d "${DEST}/windows-amd64/lib" ]; then
  echo "already present: ${DEST}/windows-amd64 (delete it to re-download)"
  exit 0
fi

for tool in curl tar sha256sum; do
  command -v "${tool}" >/dev/null || { echo "missing tool: ${tool}" >&2; exit 1; }
done

TMP="$(mktemp -d)"
trap 'rm -rf "${TMP}"' EXIT

echo "downloading ${URL}"
curl -fL --retry 3 --retry-delay 2 -o "${TMP}/${ASSET}" "${URL}"

echo "verifying checksum"
echo "${SHA256}  ${TMP}/${ASSET}" | sha256sum -c -

mkdir -p "${DEST}"
if command -v zstd >/dev/null; then
  tar --use-compress-program=zstd -xf "${TMP}/${ASSET}" -C "${DEST}"
else
  tar -xf "${TMP}/${ASSET}" -C "${DEST}"
fi

test -f "${DEST}/windows-amd64/lib/libLLVMCore.a" ||
  { echo "unpack failed: libLLVMCore.a missing" >&2; exit 1; }
echo "installed ${DEST}/windows-amd64"