#!/bin/bash
set -e

# Configuration
LLVM_VERSION="18"
ARCH="amd64"
TARGET_DIR="../llvm/18.1.3/$ARCH"

mkdir -p "$TARGET_DIR/lib"
mkdir -p "$TARGET_DIR/include"

echo "Copying headers..."
rm -rf "$TARGET_DIR/include/llvm-c" "$TARGET_DIR/include/llvm"
cp -r /usr/lib/llvm-$LLVM_VERSION/include/llvm-c "$TARGET_DIR/include/"
cp -r /usr/lib/llvm-$LLVM_VERSION/include/llvm "$TARGET_DIR/include/"
# lld/Common is only needed by lld_link.cpp; some distro packages ship it
# under llvm-$LLVM_VERSION/include, others under /usr/include.
for lld_inc in "/usr/lib/llvm-$LLVM_VERSION/include/lld" "/usr/include/lld"; do
    if [[ -d "$lld_inc" ]]; then
        rm -rf "$TARGET_DIR/include/lld"
        cp -r "$lld_inc" "$TARGET_DIR/include/"
        break
    fi
done

echo "Copying static libraries..."
for lib in $(llvm-config-$LLVM_VERSION --libs --link-static | xargs -n1); do
    libname="${lib#-l}"
    # Exclude Polly as it might be missing
    if [[ "$libname" != "Polly" && "$libname" != "PollyISL" ]]; then
        cp "/usr/lib/llvm-$LLVM_VERSION/lib/lib${libname}.a" "$TARGET_DIR/lib/"
    fi
done

# Dynamic library is not committed to repo to avoid large file size.
# Users should install LLVM system-wide for dynamic linking.
# echo "Copying dynamic library..."
# cp "/usr/lib/llvm-$LLVM_VERSION/lib/libLLVM-$LLVM_VERSION.so" "$TARGET_DIR/lib/libLLVM.so"

echo "Done."
