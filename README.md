# go-llvm

A Go wrapper for LLVM C API, designed for portability and static linking support.
This library serves as a replacement for `tinygo.org/x/go-llvm` with enhanced support for self-contained builds (embedding LLVM static libraries) and specific compatibility with `yaklang` ecosystem.

## Features

*   **Static Linking**: Embeds LLVM static libraries (`.a`) to allow building binaries that do not depend on system LLVM installation at runtime (though system libraries like `libz`, `libzstd`, `libffi` are still required).
*   **Dynamic Linking**: Supports linking against system `libLLVM.so` for faster development builds.
*   **JIT Support**: Provides `ExecutionEngine` support, specifically using the LLVM Interpreter mode to ensure compatibility with `RunFunction` and `GenericValue` argument passing on modern LLVM versions (where MCJIT has dropped support for some features).
*   **API Compatibility**: Interfaces are designed to be largely compatible with `tinygo.org/x/go-llvm`.

## Requirements

*   **Supported Platforms**:
    *   **Linux (amd64)**: Primary support.
    *   **macOS (amd64/arm64)**: Supported via Homebrew `llvm@18`.
    *   **Windows (amd64)**: Supported via MSYS2 MinGW64.
*   **LLVM 18**: Libraries are generated from LLVM 18.1.3.
*   **System Dependencies** (for static linking):
    *   **Linux**: `libz` (`zlib1g-dev`), `libzstd` (`libzstd-dev`), `libffi` (`libffi-dev`), `libxml2` (`libxml2-dev`), `libstdc++`.
    *   **macOS**: `zlib`, `libxml2` (usually provided by system/SDK), `libcompression`.
    *   **Windows (MSYS2)**: `mingw-w64-x86_64-zlib`, `mingw-w64-x86_64-zstd`, `mingw-w64-x86_64-libffi`, `mingw-w64-x86_64-libxml2`.

## Installation

```bash
go get github.com/yaklang/go-llvm
```

## Usage

### Default (Static Linking)

By default, `go-llvm` uses static linking. This requires the static libraries to be present in `llvm/18.1.3/<ARCH>/lib`.

```bash
go build .
```

### Dynamic Linking

To link against system `libLLVM` (requires LLVM installed and in library path):

```bash
go build -tags llvm_dynamic .
```

### Example: JIT Compilation

(See `examples/simple_jit.go`)

## Development & Library Generation

The `llvm/` directory contains platform-specific static libraries. You must populate these for your platform if you are developing on a new machine or cross-compiling.

### Linux (Debian/Ubuntu)

```bash
sudo apt install llvm-18-dev libz-dev libzstd-dev libffi-dev libxml2-dev
cd scripts
./generate_libs.sh
```

### macOS

Requires Homebrew:

```bash
brew install llvm@18
cd scripts
./generate_libs_macos.sh
```

### Windows (MSYS2 MinGW64)

Requires MSYS2 with MinGW64 environment:

```bash
pacman -S mingw-w64-x86_64-llvm mingw-w64-x86_64-clang
cd scripts
./generate_libs_windows.sh
```

### Updating Linker Flags

The `llvm_config_*` files contain `#cgo` directives. If you change LLVM version or dependencies, you may need to update these files or run `scripts/gen_platform_configs.sh` (on Linux) to regenerate the base list.

*   `llvm_config_linux_static.go` / `_dynamic.go`
*   `llvm_config_darwin_static.go` / `_dynamic.go`
*   `llvm_config_windows_static.go` / `_dynamic.go`

## License

Same as LLVM (Apache 2.0 License with LLVM exceptions).
