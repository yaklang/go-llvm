#!/bin/bash
set -e

# Configuration
# Assuming Homebrew installation
LLVM_PREFIX="$(brew --prefix llvm@18)"
ARCH="$(uname -m)"
TARGET_DIR="../llvm/18.1.3/$ARCH"

if [ ! -d "$LLVM_PREFIX" ]; then
    echo "Error: llvm@18 not found. Please run: brew install llvm@18"
    exit 1
fi

mkdir -p "$TARGET_DIR/lib"
mkdir -p "$TARGET_DIR/include"

echo "Using LLVM at: $LLVM_PREFIX"

echo "Copying headers..."
# macOS uses cp -R
cp -R "$LLVM_PREFIX/include/llvm-c" "$TARGET_DIR/include/"

echo "Copying static libraries..."
# On macOS, static libs are .a
# llvm-config output might vary, but usually they are in lib/
# We filter out dynamic libraries if mixed
for lib in "$LLVM_PREFIX"/lib/*.a; do
    cp "$lib" "$TARGET_DIR/lib/"
done

echo "Copying dynamic library..."
# macOS uses .dylib
if [ -f "$LLVM_PREFIX/lib/libLLVM.dylib" ]; then
    cp "$LLVM_PREFIX/lib/libLLVM.dylib" "$TARGET_DIR/lib/libLLVM.dylib"
else
    echo "Warning: libLLVM.dylib not found (optional for static build)"
fi

echo "Done. Libraries copied to $TARGET_DIR"
