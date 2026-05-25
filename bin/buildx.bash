#!/usr/bin/env bash
set -e

# see https://lucor.dev/post/cross-compile-golang-fyne-project-using-zig/


export CMD=gitall
export BUILD_DIR=build/github.com/jkassis

# Zig's 'x86_64' is Go's 'amd64'

# darwin/arm64
export MACOS_MIN_VER=13.3
export MACOS_SDK_PATH="/Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk/"
mkdir -p dist/darwin-arm64
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=arm64
# export CGO_LDFLAGS="-target arm64"
# export CGO_LDFLAGS="-mmacosx-version-min=${MACOS_MIN_VER} -target arm64 --sysroot ${MACOS_SDK_PATH} -F/System/Library/Frameworks -L/usr/lib" \
# export CC="zig cc -target aarch64-macos"
# export CXX="zig c++ -target aarch64-macos"
export CC="zig cc -target aarch64-linux"
export CCX="zig c++ -target aarch64-linux"
# export CC="zig cc -target aarch64-macos-gnu -I<sysroot>/usr/include -F<sysroot>/System/Library/Frameworks -L<sysroot>/usr/lib"
# export CXX="zig c++ -target aarch64-macos-gnu -I<sysroot>/usr/include -F<sysroot>/System/Library/Frameworks -L<sysroot>/usr/lib"
go build -o ${BUILD_DIR}/${CMD}-darwin-${MACOS_MIN_VER}-arm64 ./cmd/...


# -I<sysroot>/usr/include, -F<sysroot>/System/Library/Frameworks and -L<sysroot>/usr/lib

# darwin/amd64
export MACOS_MIN_VER=13.3
export MACOS_SDK_PATH="/Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk/"
mkdir -p dist/darwin-amd64
CGO_ENABLED=1 \
GOOS=darwin \
GOARCH=amd64 \
CGO_LDFLAGS="-mmacosx-version-min=${MACOS_MIN_VER} --sysroot ${MACOS_SDK_PATH} -F/System/Library/Frameworks -L/usr/lib" \
CC="zig cc -target x86_64-macos -iwithsysroot /usr/include -iframeworkwithsysroot /System/Library/Frameworks" \
CXX="zig c++ -target x86_64-macos -iwithsysroot /usr/include -iframeworkwithsysroot /System/Library/Frameworks" \
go build -trimpath -buildmode=pie -o ${BUILD_DIR}/${CMD}-darwin-${MACOS_MIN_VER}-amd64  ./cmd/...


# darwin/amd64
export MACOS_MIN_VER=10.14
export MACOS_SDK_PATH="/path/to/macOS/sdk"
mkdir -p dist/darwin-amd64
CGO_ENABLED=1 \
GOOS=darwin \
GOARCH=amd64 \
CGO_LDFLAGS="-mmacosx-version-min=${MACOS_MIN_VER} --sysroot ${MACOS_SDK_PATH} -F/System/Library/Frameworks -L/usr/lib" \
CC="zig cc -mmacosx-version-min=${MACOS_MIN_VER} -target x86_64 -isysroot ${MACOS_SDK_PATH} -iwithsysroot /usr/include -iframeworkwithsysroot /System/Library/Frameworks" \
CXX="zig c++ -mmacosx-version-min=${MACOS_MIN_VER} -target x86_64 -isysroot ${MACOS_SDK_PATH} -iwithsysroot /usr/include -iframeworkwithsysroot /System/Library/Frameworks" \
go build -trimpath -buildmode=pie -o ${BUILD_DIR}/${CMD}-darwin-${MACOS_MIN_VER}-amd64  ./cmd/...


# darwin/arm64
export MACOS_MIN_VER=11.1
export MACOS_SDK_PATH="/path/to/macOS/sdk"
mkdir -p dist/darwin-arm64
CGO_ENABLED=1 \
GOOS=darwin \
GOARCH=amd64 \
CGO_LDFLAGS="-mmacosx-version-min=${MACOS_MIN_VER} --sysroot ${MACOS_SDK_PATH} -F/System/Library/Frameworks -L/usr/lib" \
CC="zig cc -mmacosx-version-min=${MACOS_MIN_VER} -target aarch64-macos-gnu -isysroot ${MACOS_SDK_PATH} -iwithsysroot /usr/include -iframeworkwithsysroot /System/Library/Frameworks" \
CXX="zig c++ -mmacosx-version-min=${MACOS_MIN_VER} -target aarch64-macos-gnu -isysroot ${MACOS_SDK_PATH} -iwithsysroot /usr/include -iframeworkwithsysroot /System/Library/Frameworks" \
go build -ldflags "-s -w" -buildmode=pie -trimpath -o ${BUILD_DIR}/${CMD}-darwin-${MACOS_MIN_VER}-arm64 ./cmd/...


# linux/amd64
mkdir -p dist/linux-amd64
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=amd64 \
CC="zig cc -target x86_64-linux -isystem /usr/include -L/usr/lib/x86_64-linux-gnu"  \
CXX="zig c++ -target x86_64-linux -isystem /usr/include -L/usr/lib/x86_64-linux-gnu" \
go build -trimpath -o ${BUILD_DIR}/${CMD}-linux-amd64 ./cmd/...


# linux/arm64
mkdir -p dist/linux-arm64
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=arm64 \
PKG_CONFIG_LIBDIR=/usr/lib/aarch64-linux-gnu/pkgconfig \
CC="zig cc -target aarch64-linux-gnu -isystem /usr/include -L/usr/lib/aarch64-linux-gnu"  \
CXX="zig c++ -target aarch64-linux-gnu -isystem /usr/include -L/usr/lib/aarch64-linux-gnu" \
go build -trimpath -o ${BUILD_DIR}/${CMD}-linux-arm64 ./cmd/...


# windows/amd64
mkdir -p dist/windows-amd64
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC="zig cc -target x86_64-windows-gnu" \
CXX="zig c++ -target x86_64-windows-gnu" \
go build -trimpath -ldflags='-H=windowsgui' -o ${BUILD_DIR}/${CMD}-windows-amd64 ./cmd/...
