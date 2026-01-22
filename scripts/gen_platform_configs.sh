#!/bin/bash

# Extract LLVM libs from linux config
LIBS=$(grep "LDFLAGS:" ../llvm_config_linux_static.go | sed 's/.*-lstdc++//' | sed 's/ -lffi//')

# Generate Darwin Static
cat <<EOF > ../llvm_config_darwin_static.go
//go:build !llvm_dynamic && darwin
// +build !llvm_dynamic,darwin

package llvm

/*
#cgo darwin,amd64 CFLAGS: -I\${SRCDIR}/llvm/18.1.3/amd64/include
#cgo darwin,arm64 CFLAGS: -I\${SRCDIR}/llvm/18.1.3/arm64/include
#cgo darwin,amd64 LDFLAGS: -L\${SRCDIR}/llvm/18.1.3/amd64/lib -lc++ -lz -lcurses -lm -lxml2 -lcompression $LIBS
#cgo darwin,arm64 LDFLAGS: -L\${SRCDIR}/llvm/18.1.3/arm64/lib -lc++ -lz -lcurses -lm -lxml2 -lcompression $LIBS
*/
import "C"
EOF

# Generate Darwin Dynamic
cat <<EOF > ../llvm_config_darwin_dynamic.go
//go:build llvm_dynamic && darwin
// +build llvm_dynamic,darwin

package llvm

/*
#cgo darwin,amd64 CFLAGS: -I\${SRCDIR}/llvm/18.1.3/amd64/include
#cgo darwin,arm64 CFLAGS: -I\${SRCDIR}/llvm/18.1.3/arm64/include
#cgo darwin,amd64 LDFLAGS: -L\${SRCDIR}/llvm/18.1.3/amd64/lib -lLLVM
#cgo darwin,arm64 LDFLAGS: -L\${SRCDIR}/llvm/18.1.3/arm64/lib -lLLVM
*/
import "C"
EOF

# Generate Windows Static
cat <<EOF > ../llvm_config_windows_static.go
//go:build !llvm_dynamic && windows
// +build !llvm_dynamic,windows

package llvm

/*
#cgo windows,amd64 CFLAGS: -I\${SRCDIR}/llvm/18.1.3/amd64/include
#cgo windows,amd64 LDFLAGS: -L\${SRCDIR}/llvm/18.1.3/amd64/lib -lstdc++ -static -lpthread -lws2_32 -lole32 -luuid -lkernel32 -luser32 -lgdi32 -lwinspool -lshell32 -ladvapi32 -lversion $LIBS
*/
import "C"
EOF

# Generate Windows Dynamic
cat <<EOF > ../llvm_config_windows_dynamic.go
//go:build llvm_dynamic && windows
// +build llvm_dynamic,windows

package llvm

/*
#cgo windows,amd64 CFLAGS: -I\${SRCDIR}/llvm/18.1.3/amd64/include
#cgo windows,amd64 LDFLAGS: -L\${SRCDIR}/llvm/18.1.3/amd64/lib -lLLVM
*/
import "C"
EOF

echo "Generated platform configs."
