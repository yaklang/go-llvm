#!/usr/bin/env bash
set -euo pipefail

# Cross-build LLVM 18.1.3 + LLD for Windows amd64 (MinGW) on a Linux host and
# populate llvm/18.1.3/windows-amd64.
#
# Why: MSYS2 no longer ships LLVM 18 and no public download provides static
# MinGW archives for it, so the Windows tree has to be built from source. This
# script only needs a MinGW-w64 cross toolchain and clang; it does not need a
# Windows machine.
#
# Requirements (Debian/Ubuntu):
#   sudo apt install clang cmake ninja-build g++-mingw-w64-x86-64
#
# Usage:
#   scripts/build_windows_llvm18.sh [work-dir]
#
# work-dir defaults to ../llvm18-windows-build next to this checkout. The build
# takes roughly 10-20 minutes on 20 cores and about 8 GB of disk.

LLVM_VERSION=18.1.3
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
WORK="${1:-${REPO_ROOT}/../llvm18-windows-build}"
SRC="${WORK}/llvm-project-${LLVM_VERSION}.src"
TARGET="${REPO_ROOT}/llvm/18.1.3/windows-amd64"

for tool in clang clang++ cmake ninja x86_64-w64-mingw32-gcc x86_64-w64-mingw32-g++; do
  if ! command -v "${tool}" >/dev/null 2>&1; then
    echo "build_windows_llvm18: missing ${tool}" >&2
    exit 1
  fi
done

mkdir -p "${WORK}"

if [[ ! -d "${SRC}" ]]; then
  archive="${WORK}/llvm-project-${LLVM_VERSION}.src.tar.xz"
  url="https://github.com/llvm/llvm-project/releases/download/llvmorg-${LLVM_VERSION}/llvm-project-${LLVM_VERSION}.src.tar.xz"
  echo "[windows] downloading ${url}"
  curl -fL -o "${archive}" "${url}"
  tar -C "${WORK}" -xf "${archive}"
fi

common_cmake=(
  -G Ninja
  -S "${SRC}/llvm"
  -DCMAKE_BUILD_TYPE=Release
  -DLLVM_TARGETS_TO_BUILD=X86
  -DCMAKE_CXX_FLAGS=-include\ cstdint
  -DLLVM_INCLUDE_TESTS=OFF
  -DLLVM_INCLUDE_DOCS=OFF
  -DLLVM_INCLUDE_EXAMPLES=OFF
  -DLLVM_INCLUDE_BENCHMARKS=OFF
  -DLLVM_ENABLE_ZLIB=OFF
  -DLLVM_ENABLE_ZSTD=OFF
  -DLLVM_ENABLE_LIBXML2=OFF
  -DLLVM_ENABLE_TERMINFO=OFF
  -DLLVM_ENABLE_LIBEDIT=OFF
  -DLLVM_ENABLE_BINDINGS=OFF
  -DLLVM_ENABLE_PLUGINS=OFF
)

if [[ ! -x "${WORK}/build-native/bin/llvm-tblgen" ]]; then
  echo "[windows] building host llvm-tblgen"
  cmake "${common_cmake[@]}" -B "${WORK}/build-native" -DLLVM_BUILD_TOOLS=OFF
  ninja -C "${WORK}/build-native" llvm-tblgen
fi

echo "[windows] configuring the MinGW cross build"
cmake "${common_cmake[@]}" -B "${WORK}/build-win" \
  -DCMAKE_SYSTEM_NAME=Windows \
  -DCMAKE_C_COMPILER=clang \
  -DCMAKE_C_COMPILER_TARGET=x86_64-w64-windows-gnu \
  -DCMAKE_CXX_COMPILER=clang++ \
  -DCMAKE_CXX_COMPILER_TARGET=x86_64-w64-windows-gnu \
  -DCMAKE_ASM_COMPILER=clang \
  -DCMAKE_ASM_COMPILER_TARGET=x86_64-w64-windows-gnu \
  -DCMAKE_RC_COMPILER=x86_64-w64-mingw32-windres \
  -DCMAKE_FIND_ROOT_PATH=/usr/x86_64-w64-mingw32 \
  -DCMAKE_FIND_ROOT_PATH_MODE_PROGRAM=NEVER \
  -DCMAKE_FIND_ROOT_PATH_MODE_LIBRARY=BOTH \
  -DCMAKE_FIND_ROOT_PATH_MODE_INCLUDE=BOTH \
  -DCMAKE_CXX_STANDARD_LIBRARIES=-lwinpthread \
  -DLLVM_ENABLE_PROJECTS=lld \
  -DLLVM_NATIVE_TOOL_DIR="${WORK}/build-native/bin" \
  -DLLVM_BUILD_TOOLS=OFF \
  -DLLVM_BUILD_UTILS=OFF \
  -DLLVM_HOST_TRIPLE=x86_64-w64-mingw32 \
  -DLLVM_DEFAULT_TARGET_TRIPLE=x86_64-w64-mingw32

echo "[windows] building"
ninja -C "${WORK}/build-win" -j"$(nproc)"

echo "[windows] populating ${TARGET}"
if [[ -d "${TARGET}" ]]; then
  find "${TARGET}" -mindepth 1 -delete
fi
mkdir -p "${TARGET}/include" "${TARGET}/lib"
cp -R "${SRC}/llvm/include/llvm" "${SRC}/llvm/include/llvm-c" "${TARGET}/include/"
cp -R "${SRC}/lld/include/lld" "${TARGET}/include/"
cp "${WORK}/build-win/include/llvm/Config/"* "${TARGET}/include/llvm/Config/"
cp "${WORK}"/build-win/lib/libLLVM*.a "${WORK}"/build-win/lib/liblld*.a "${TARGET}/lib/"
x86_64-w64-mingw32-ranlib "${TARGET}"/lib/*.a

for lib in libLLVMCore.a libLLVMSupport.a liblldCommon.a liblldMinGW.a; do
  if [[ ! -f "${TARGET}/lib/${lib}" ]]; then
    echo "build_windows_llvm18: ${lib} missing from the build" >&2
    exit 1
  fi
done

echo "[windows] done: ${TARGET}"