// LLVM's LTO backend (libLLVMLTO.a, LTOBackend.cpp.o) references
// getPollyPluginInfo(), which is defined by the optional Polly plugin. Polly is
// intentionally not bundled (see scripts/generate_libs.sh). Without lld this
// object was never pulled into the link, so the reference stayed inert; lld's
// liblldELF.a pulls LLVM LTO symbols, activating LTOBackend.cpp.o and turning
// this into a strong undefined reference. Provide a NULL stub so the static
// link resolves. At runtime the LTO backend treats a null return as "no Polly
// plugin", which is exactly the intended behavior when Polly is absent.
//
// The mangled symbol is _Z18getPollyPluginInfov (getPollyPluginInfo, no args);
// C++ return types are not mangled, so a void* stub matches the reference.
extern "C++" void *getPollyPluginInfo() { return nullptr; }