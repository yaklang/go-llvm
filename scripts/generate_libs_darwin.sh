#!/usr/bin/env bash
set -euo pipefail

# Populate llvm/18.1.3/darwin-<arch> from Homebrew's llvm@18.
# The generated directory is local to this checkout; the repository only ships
# the Linux archives.
#
#   brew install llvm@18
#   scripts/generate_libs_darwin.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

case "$(uname -m)" in
  x86_64 | amd64) GOARCH=amd64 ;;
  arm64 | aarch64) GOARCH=arm64 ;;
  *)
    echo "generate_libs_darwin: unsupported machine $(uname -m)" >&2
    exit 1
    ;;
esac

PREFIX="${LLVM_PREFIX:-}"
if [[ -z "${PREFIX}" ]]; then
  if ! command -v brew >/dev/null 2>&1; then
    echo "generate_libs_darwin: set LLVM_PREFIX or install Homebrew" >&2
    exit 1
  fi
  PREFIX="$(brew --prefix llvm@18)"
fi
if [[ ! -d "${PREFIX}/include/llvm-c" ]]; then
  echo "generate_libs_darwin: ${PREFIX} does not look like an LLVM install" >&2
  exit 1
fi

LIB_TARGET="${REPO_ROOT}/native/darwin-${GOARCH}/lib"
INCLUDE_TARGET="${REPO_ROOT}/llvm/18.1.3/include"
mkdir -p "${INCLUDE_TARGET}" "${LIB_TARGET}"

echo "[darwin] copying headers from ${PREFIX}/include"
rm -rf "${INCLUDE_TARGET}/llvm-c" "${INCLUDE_TARGET}/llvm" "${INCLUDE_TARGET}/lld"
cp -R "${PREFIX}/include/llvm-c" "${INCLUDE_TARGET}/"
cp -R "${PREFIX}/include/llvm" "${INCLUDE_TARGET}/"
if [[ -d "${PREFIX}/include/lld" ]]; then
  cp -R "${PREFIX}/include/lld" "${INCLUDE_TARGET}/"
else
  echo "[darwin] warning: ${PREFIX}/include/lld missing; Mach-O lld headers unavailable" >&2
fi

echo "[darwin] copying static archives from ${PREFIX}/lib"
shopt -s nullglob
libs=("${PREFIX}"/lib/libLLVM*.a "${PREFIX}"/lib/liblld*.a)
shopt -u nullglob
if [[ ${#libs[@]} -eq 0 ]]; then
  echo "generate_libs_darwin: no static archives under ${PREFIX}/lib" >&2
  exit 1
fi
cp "${libs[@]}" "${LIB_TARGET}/"

echo "[darwin] done: ${LIB_TARGET}"
echo "[darwin] run scripts/trim_platform_assets.py to drop non-host backends and dedupe headers"