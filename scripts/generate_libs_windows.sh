#!/bin/bash
set -e

# Configuration
# Assuming MSYS2 / MinGW64 environment
# Adjust LLVM_PREFIX if installed elsewhere
LLVM_PREFIX="/mingw64"
ARCH="amd64"
TARGET_DIR="../llvm/18.1.3/$ARCH"

if [ ! -f "$LLVM_PREFIX/bin/llvm-config.exe" ]; then
    echo "Error: llvm-config not found in /mingw64/bin. Are you in MSYS2 MinGW64 shell?"
    echo "Please install: pacman -S mingw-w64-x86_64-llvm"
    exit 1
fi

mkdir -p "$TARGET_DIR/lib"
mkdir -p "$TARGET_DIR/include"

echo "Copying headers..."
cp -r "$LLVM_PREFIX/include/llvm-c" "$TARGET_DIR/include/"

echo "Copying static libraries..."
# MinGW libs are usually .a
# llvm-config --libs gives -lLLVM..., we need to find files
# We can just copy all .a files from lib that match libLLVM*.a
find "$LLVM_PREFIX/lib" -maxdepth 1 -name "libLLVM*.a" -exec cp {} "$TARGET_DIR/lib/" \;

echo "Copying dynamic library..."
# On Windows, DLLs are in bin/, import libs (.dll.a) in lib/
if [ -f "$LLVM_PREFIX/bin/libLLVM.dll" ]; then
    cp "$LLVM_PREFIX/bin/libLLVM.dll" "$TARGET_DIR/lib/"
fi
# Copy import lib if exists
if [ -f "$LLVM_PREFIX/lib/libLLVM.dll.a" ]; then
    cp "$LLVM_PREFIX/lib/libLLVM.dll.a" "$TARGET_DIR/lib/"
fi

echo "Done."
