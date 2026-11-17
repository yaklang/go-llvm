#!/usr/bin/env bash
set -euo pipefail

# Populate llvm/18.1.3/linux-<arch> from a local LLVM 18 installation.
# Only needed to regenerate the committed Linux archives or to build on a
# non-amd64 Linux host. Normal builds use the committed linux-amd64 files.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

case "$(uname -m)" in
  x86_64 | amd64) GOARCH=amd64 ;;
  aarch64 | arm64) GOARCH=arm64 ;;
  *)
    echo "generate_libs_linux: unsupported machine $(uname -m)" >&2
    exit 1
    ;;
esac

LLVM_CONFIG="${LLVM_CONFIG:-llvm-config-18}"
if ! command -v "${LLVM_CONFIG}" >/dev/null 2>&1; then
  echo "generate_libs_linux: ${LLVM_CONFIG} not found; install llvm-18-dev and liblld-18-dev" >&2
  exit 1
fi

PREFIX="$("${LLVM_CONFIG}" --prefix)"
LIB_TARGET="${REPO_ROOT}/native/linux-${GOARCH}/lib"
INCLUDE_TARGET="${REPO_ROOT}/llvm/18.1.3/linux-${GOARCH}/include"
mkdir -p "${INCLUDE_TARGET}" "${LIB_TARGET}"

echo "[linux] copying headers from ${PREFIX}/include"
rm -rf "${INCLUDE_TARGET}/llvm-c" "${INCLUDE_TARGET}/llvm" "${INCLUDE_TARGET}/lld"
cp -R "${PREFIX}/include/llvm-c" "${INCLUDE_TARGET}/"
cp -R "${PREFIX}/include/llvm" "${INCLUDE_TARGET}/"
if [[ -d "${PREFIX}/include/lld" ]]; then
  cp -R "${PREFIX}/include/lld" "${INCLUDE_TARGET}/"
else
  echo "[linux] warning: ${PREFIX}/include/lld missing; install liblld-18-dev" >&2
fi

echo "[linux] copying static archives from ${PREFIX}/lib"
shopt -s nullglob
libs=("${PREFIX}"/lib/libLLVM*.a "${PREFIX}"/lib/liblld*.a)
shopt -u nullglob
if [[ ${#libs[@]} -eq 0 ]]; then
  echo "generate_libs_linux: no static archives under ${PREFIX}/lib" >&2
  exit 1
fi
cp "${libs[@]}" "${LIB_TARGET}/"

echo "[linux] done: ${LIB_TARGET}"