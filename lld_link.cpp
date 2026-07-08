// lld_link.cpp — C++ shim that drives lld's in-process ELF linker via its C++
// driver API (lld has no C API). Exposes a C ABI consumed by lld.go via cgo.
//
// Design notes:
//   * exitEarly=false: lld returns (bool) instead of calling exit() after a
//     completed link. Normal link failures come back as `false` with
//     diagnostics written to the stderr stream; the caller surfaces them.
//   * disableOutput=false: the output file is actually written.
//   * A residual risk remains: a catastrophic lld/LLVM fatal error can call
//     abort() regardless of exitEarly. ssa2llvm treats that as a hard compiler
//     crash (one link per process; output is in a temp work dir).
//   * lld registers process-global signal handlers on first link; they persist.
//     This is the accepted "bounded state coupling" of lld-as-cgo-library.
#include <lld/Common/Driver.h>

#include <llvm/ADT/ArrayRef.h>
#include <llvm/Support/raw_ostream.h>

#include <cstdlib>
#include <cstring>
#include <string>

// Declare the ELF driver's link() function. LLVM 18's lld/Common/Driver.h
// exposes the per-driver link() entry points only via this macro (it expands
// to a namespace-scope declaration of lld::elf::link). The definition is in
// liblldELF.a (bundled). We call the ELF driver directly so no -flavor parsing
// is needed; exitEarly=false makes lld return instead of calling exit().
LLD_HAS_DRIVER(elf)

extern "C" {

// Drive lld's ELF linker in-process.
//   argv: argc NUL-terminated C strings (argv[0] is a dummy program name).
//   stdout_buf/stderr_buf: out params, malloc'd on return; free with
//     go_llvm_lld_free. May be NULL/0 when empty.
// Returns 0 on success, 1 on link failure.
int go_llvm_lld_link(int argc, char **argv,
                     char **stdout_buf, int *stdout_len,
                     char **stderr_buf, int *stderr_len) {
  *stdout_buf = nullptr;
  *stdout_len = 0;
  *stderr_buf = nullptr;
  *stderr_len = 0;

  std::string out_str;
  std::string err_str;
  llvm::raw_string_ostream out_os(out_str);
  llvm::raw_string_ostream err_os(err_str);

  llvm::ArrayRef<const char *> args(const_cast<const char **>(argv),
                                    static_cast<size_t>(argc));
  bool ok = lld::elf::link(args, out_os, err_os,
                           /*exitEarly=*/false, /*disableOutput=*/false);

  out_os.flush();
  err_os.flush();

  auto copy = [](const std::string &s, char **buf, int *len) {
    *len = static_cast<int>(s.size());
    if (*len > 0) {
      *buf = static_cast<char *>(std::malloc(static_cast<size_t>(*len)));
      if (*buf) {
        std::memcpy(*buf, s.data(), static_cast<size_t>(*len));
      } else {
        *len = 0;
      }
    }
  };
  copy(out_str, stdout_buf, stdout_len);
  copy(err_str, stderr_buf, stderr_len);

  return ok ? 0 : 1;
}

void go_llvm_lld_free(char *buf) {
  if (buf) {
    std::free(buf);
  }
}

} // extern "C"