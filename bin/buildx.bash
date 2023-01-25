#!/usr/bin/env bash
set -e # exit when any command fails
mkdir -p build

# This is complicated...
# currently using Docker image built from github.com/jkassis/xgo
# might want to try https://github.com/crazy-max/goxx
# or https://github.com/techknowlogick/xgo

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

# do the builds
TARGETS="linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64"
docker run --rm \
    -v "$PWD"/build:/build \
    -v "$GOPATH"/xgo-cache:/deps-cache:ro \
    -v "$PWD":/source \
    -e FLAG_V=false \
    -e FLAG_X=false \
    -e FLAG_RACE=false \
    -e FLAG_LDFLAGS="-w -s" \
    -e FLAG_BUILDMODE=default \
    -e TARGETS="$TARGETS" \
    jkassis/xgo:1.19.5 ./cmd/

# make them executable
chmod 555 build/github.com/jkassis/*
