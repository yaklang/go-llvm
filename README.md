# go-llvm

A Go wrapper for LLVM C API, designed for portability and static linking support.
This library serves as a replacement for `tinygo.org/x/go-llvm` with enhanced support for self-contained builds (embedding LLVM static libraries) and specific compatibility with `yaklang` ecosystem.

## Features

*   **Static Linking**: Embeds the LLVM 18.1.3 C/C++ headers and static libraries (`.a`) so builds do not use any system LLVM install (system libraries like `libz`, `libzstd`, `libffi` are still required).
*   **Per-platform layout**: Assets live under `llvm/18.1.3/<goos>-<goarch>/`. `linux-amd64` and `windows-amd64` are committed; macOS archives are generated locally by the script below.
*   **Host lld driver**: the in-process linker uses `lld::elf` on Linux, `lld::macho` on macOS, and `lld::mingw` on Windows.
*   **Dynamic Linking**: Supports linking against system `libLLVM.so` for faster development builds.
*   **JIT Support**: Provides `ExecutionEngine` support, specifically using the LLVM Interpreter mode to ensure compatibility with `RunFunction` and `GenericValue` argument passing on modern LLVM versions (where MCJIT has dropped support for some features).
*   **API Compatibility**: Interfaces are designed to be largely compatible with `tinygo.org/x/go-llvm`.

## Requirements

*   **Supported Platforms**:
    *   **Linux (amd64)**: committed archives; primary CI target.
    *   **macOS (arm64/amd64)**: generate archives with Homebrew `llvm@18`.
    *   **Windows (amd64, MinGW)**: committed archives (X86 only). Regenerate with `scripts/build_windows_llvm18.sh` on Linux, or unpack the release asset with `scripts/fetch_windows_llvm18.sh`. MSYS2 no longer ships LLVM 18, so there is no distro package to copy.
*   **LLVM 18**: Headers and libraries are embedded per platform. No system LLVM installation is used for builds once the platform directory exists.
*   **System Dependencies** (for static linking):
    *   **Linux**: `libz` (`zlib1g-dev`), `libzstd` (`libzstd-dev`), `libffi` (`libffi-dev`), `libxml2` (`libxml2-dev`), `libstdc++`.
    *   **macOS**: `zlib`, `libxml2` (usually provided by system/SDK), `libcompression`.
    *   **Windows**: MinGW-w64 runtime only; the cross build disables zlib, zstd, libxml2 and terminfo, so no extra DLLs are needed beyond `libstdc++-6.dll` and `libwinpthread-1.dll` (or link them statically).

## Installation

```bash
go get github.com/yaklang/go-llvm
```

## Usage

### Default (Static Linking)

By default, `go-llvm` uses static linking. This requires the static libraries to be present in `llvm/18.1.3/<goos>-<goarch>/lib`.

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

The `llvm/` directory contains per-platform headers and static libraries under
`llvm/18.1.3/<goos>-<goarch>/`. Only `linux-amd64` is committed so `go get`
works out of the box on Linux. The scripts write the current host's tree; use a
local checkout (and a `replace` directive) to build on macOS or Windows.

### Linux (Debian/Ubuntu)

```bash
sudo apt install llvm-18-dev liblld-18-dev libz-dev libzstd-dev libffi-dev libxml2-dev
scripts/generate_libs_linux.sh
```

The committed `linux-amd64` tree keeps only the X86 backend: 172 MB of archives
instead of 296 MB. `targets_x86_only.cpp` supplies empty init symbols for the
other backends, because lld's Driver.cpp.o inlines `llvm::InitializeAllTarget*`
and would otherwise pull every backend into the link. The script above still
generates the full multi-target set; delete the non-X86 archives and rerun the
config step if you want the trimmed layout.

### macOS

```bash
brew install llvm@18
scripts/generate_libs_darwin.sh
```

### Windows (amd64, MinGW)

`llvm/18.1.3/windows-amd64` is committed: LLVM 18.1.3 + LLD cross-built for
`x86_64-w64-windows-gnu`, X86 target only, about 188 MB. `go get` works without
any extra step; MinGW-w64 provides the C compiler at build time.

Regenerate the tree on a Linux host when the LLVM version or build options
change:

```bash
sudo apt install clang cmake ninja-build g++-mingw-w64-x86-64
scripts/build_windows_llvm18.sh
```

If a checkout is missing the tree, `scripts/fetch_windows_llvm18.sh` unpacks the
same tree from the release asset after verifying its sha256.

### Updating Linker Flags

The `llvm_config_*` files contain `#cgo` directives. Update them by hand when
changing LLVM version or a platform's library list.

*   `llvm_config_linux_static.go` / `_dynamic.go`
*   `llvm_config_darwin_static.go` / `_dynamic.go`
*   `llvm_config_windows_static.go` / `_dynamic.go`

## License

Same as LLVM (Apache 2.0 License with LLVM exceptions).
