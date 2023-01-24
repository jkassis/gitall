#!/usr/bin/env bash

# GOOS - Target Operating System	GOARCH - Target Platform
# android	arm
# darwin	386
# darwin	amd64
# darwin	arm
# darwin	arm64
# dragonfly	amd64
# freebsd	386
# freebsd	amd64
# freebsd	arm
# linux	386
# linux	amd64
# linux	arm
# linux	arm64
# linux	ppc64
# linux	ppc64le
# linux	mips
# linux	mipsle
# linux	mips64
# linux	mips64le
# netbsd	386
# netbsd	amd64
# netbsd	arm
# openbsd	386
# openbsd	amd64
# openbsd	arm
# plan9	386
# plan9	amd64
# solaris	amd64
# windows	386
# windows	amd64

# Needed environment variables:
#   DEPS           - Optional list of C dependency packages to build
#   ARGS           - Optional arguments to pass to C dependency configure scripts
#   OUT            - Optional output prefix to override the package name
#   FLAG_V         - Optional verbosity flag to set on the Go builder
#   FLAG_X         - Optional flag to print the build progress commands
#   FLAG_RACE      - Optional race flag to set on the Go builder
#   FLAG_TAGS      - Optional tag flag to set on the Go builder
#   FLAG_LDFLAGS   - Optional ldflags flag to set on the Go builder
#   FLAG_BUILDMODE - Optional buildmode flag to set on the Go builder
#   FLAG_TRIMPATH  - Optional trimpath flag to set on the Go builder
#   TARGETS        - Comma separated list of build targets to compile for
#   GO_VERSION     - Bootstrapped version of Go to disable uncupported targets
#   EXT_GOPATH     - GOPATH elements mounted from the host filesystem

TARGETS="linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64"
docker run --rm \
    -v "$PWD"/build:/build \
    -v "$GOPATH"/xgo-cache:/deps-cache:ro \
    -v "$PWD":/source \
    # -e OUT="$OUT"\
    -e FLAG_V=false \
    -e FLAG_X=false \
    -e FLAG_RACE=false \
    -e FLAG_LDFLAGS="-w -s" \
    -e FLAG_BUILDMODE=default \
    -e TARGETS="$TARGETS" \
    mysteriumnetwork/xgo:1.18.0 ./cmd/