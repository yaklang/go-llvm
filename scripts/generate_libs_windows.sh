#!/usr/bin/env bash
set -euo pipefail

# Copy an existing MinGW LLVM 18 prefix into llvm/18.1.3/windows-amd64.
#
# MSYS2 stopped shipping LLVM 18, and its package only provides DLL import
# libraries, not the static libLLVM*.a archives go-llvm links against. Use
# scripts/build_windows_llvm18.sh (Linux cross build) unless you already built
# LLVM 18 yourself; in that case point LLVM_PREFIX at the install prefix:
#
#   LLVM_PREFIX=/c/src/llvm-install scripts/generate_libs_windows.sh
#
# The generated directory is local to this checkout; the repository only ships
# the Linux archives.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

PREFIX="${LLVM_PREFIX:-/mingw64}"
if [[ ! -d "${PREFIX}/include/llvm-c" ]]; then
  echo "generate_libs_windows: ${PREFIX} does not look like an MSYS2 MinGW64 prefix" >&2
  echo "run this inside an MSYS2 MINGW64 shell, or set LLVM_PREFIX" >&2
  exit 1
fi

GOARCH=amd64
TARGET="${REPO_ROOT}/llvm/18.1.3/windows-${GOARCH}"
mkdir -p "${TARGET}/include" "${TARGET}/lib"

echo "[windows] copying headers from ${PREFIX}/include"
rm -rf "${TARGET}/include/llvm-c" "${TARGET}/include/llvm" "${TARGET}/include/lld"
cp -R "${PREFIX}/include/llvm-c" "${TARGET}/include/"
cp -R "${PREFIX}/include/llvm" "${TARGET}/include/"
if [[ -d "${PREFIX}/include/lld" ]]; then
  cp -R "${PREFIX}/include/lld" "${TARGET}/include/"
else
  echo "[windows] warning: ${PREFIX}/include/lld missing; install mingw-w64-x86_64-lld" >&2
fi

echo "[windows] copying static archives from ${PREFIX}/lib"
shopt -s nullglob
libs=("${PREFIX}"/lib/libLLVM*.a "${PREFIX}"/lib/liblld*.a)
shopt -u nullglob
if [[ ${#libs[@]} -eq 0 ]]; then
  echo "generate_libs_windows: no static archives under ${PREFIX}/lib" >&2
  echo "the MSYS2 package normally ships libLLVM*.a; check for libLLVM*.dll.a only" >&2
  exit 1
fi
cp "${libs[@]}" "${TARGET}/lib/"

echo "[windows] done: ${TARGET}"