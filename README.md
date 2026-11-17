# go-llvm

A Go wrapper for LLVM C API, designed for portability and static linking support.
This library serves as a replacement for `tinygo.org/x/go-llvm` with enhanced support for self-contained builds (embedding LLVM static libraries) and specific compatibility with `yaklang` ecosystem.

## Features

*   **Static Linking**: Embeds the LLVM 18.1.3 C/C++ headers and static libraries (`.a`) so builds do not use any system LLVM install (system libraries like `libz`, `libzstd`, `libffi` are still required).
*   **One module per platform**: The main module holds the public API and one LLVM 18 header tree per platform under `llvm/18.1.3/<goos>-<goarch>/include`. Each platform's static archives live in their own Go module under `native/<goos>-<goarch>/`, with a `link.go` that carries the cgo link flags. `native/linux-amd64`, `native/windows-amd64` and `native/darwin-arm64` are shipped. Callers still import `github.com/yaklang/go-llvm` only.
*   **Host lld driver**: the in-process linker uses `lld::elf` on Linux, `lld::macho` on macOS, and `lld::mingw` on Windows.
*   **Dynamic Linking**: Supports linking against system `libLLVM.so` for faster development builds.
*   **JIT Support**: Provides `ExecutionEngine` support, specifically using the LLVM Interpreter mode to ensure compatibility with `RunFunction` and `GenericValue` argument passing on modern LLVM versions (where MCJIT has dropped support for some features).
*   **API Compatibility**: Interfaces are designed to be largely compatible with `tinygo.org/x/go-llvm`.

## Requirements

*   **Supported Platforms**:
    *   **Linux (amd64)**: shipped archives; primary CI target.
    *   **macOS (arm64)**: shipped archives. `darwin/amd64` is not shipped; build against a Homebrew `llvm@18` install instead.
    *   **Windows (amd64, MinGW)**: shipped archives, X86 target only. Regenerate with `scripts/build_windows_llvm18.sh` on Linux; MSYS2 no longer ships LLVM 18.
*   **LLVM 18**: Headers and libraries are embedded per platform, each tree matching its own build. No system LLVM installation is used for builds.
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

By default, `go-llvm` uses static linking. The platform archives are pulled in by
the platform module; no extra flag or local checkout is needed.

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

## Module layout

```
go-llvm                          module github.com/yaklang/go-llvm   67 MiB
  llvm/18.1.3/linux-amd64/include/     23.8 MiB, LLVM_NATIVE_ARCH X86
  llvm/18.1.3/windows-amd64/include/   19.3 MiB, LLVM_NATIVE_ARCH X86
  llvm/18.1.3/darwin-arm64/include/    23.8 MiB, LLVM_NATIVE_ARCH AArch64
  native/linux-amd64/            module .../native/linux-amd64      296 MiB
    go.mod
    link.go                      #cgo LDFLAGS for this platform
    lib/*.a
  native/windows-amd64/          module .../native/windows-amd64    164 MiB
  native/darwin-arm64/           module .../native/darwin-arm64     220 MiB
```

Modules in this repository are split because the Go module proxy rejects a
module whose unpacked source tree exceeds 500 MiB. Splitting keeps every module
far below that limit; the largest one is 296 MiB.

Headers stay in the main module because cgo include flags only apply to the
package that is being compiled. Each platform keeps its own copy: `llvm/Config/llvm-config.h`
defines `LLVM_NATIVE_ARCH`, and the inline `LLVMInitializeNativeTarget` in
`llvm-c/Target.h` registers exactly that backend. A shared header set would make
macOS arm64 register X86 and fail with "No available targets".

`go build` for one platform downloads the main module plus only that platform's
module. The main `go.mod` requires all three platform modules, so `go mod tidy`
and `go mod download all` fetch all four modules (about 750 MiB together). A
bare `go mod download` fetches three of them on Linux: main, linux-amd64 and
windows-amd64. That is a download-size cost only; no module hits the size limit.

`link.go` in each platform module holds `-L${SRCDIR}/lib` and the `-l` flags, so
the flags travel with the archives and `${SRCDIR}` resolves inside the module
cache. The main module's `llvm_config_<platform>_static.go` only adds the header
include path and blank-imports the platform module.

## Regenerating archives

```bash
# Linux host, from llvm-18-dev + liblld-18-dev
scripts/generate_libs_linux.sh

# macOS host, from Homebrew llvm@18
scripts/generate_libs_darwin.sh

# Windows archives, cross-built on Linux
scripts/build_windows_llvm18.sh
```

The scripts copy archives into `native/<goos>-<goarch>/lib` and headers into
`llvm/18.1.3/<goos>-<goarch>/include`. The link flag list in `link.go` is written
by hand; update it when the archive set changes.

## Updating Linker Flags

The cgo flags live in two places:

*   `llvm_config_{linux,windows,darwin}_static.go` — build tags, header include
    path, and the blank import of the platform module.
*   `native/<goos>-<goarch>/link.go` — `#cgo LDFLAGS` with `-L${SRCDIR}/lib` and
    the `-lLLVM*` / `-llld*` list.

`llvm_config_*_dynamic.go` keeps the `llvm_dynamic` build tag working against a
system LLVM.

## License

Same as LLVM (Apache 2.0 License with LLVM exceptions).