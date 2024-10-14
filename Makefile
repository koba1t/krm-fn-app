
# Makefile for building a Go project

# build a wasi module
# https://go.dev/blog/wasi
.PHONY: build-wasm
build-wasm:
	GOOS=wasip1 GOARCH=wasm go build -ldflags="-s -w" -trimpath -o ./bin/fnapp.wasm main.go

.PHONY: test-kustomize-wasmtime
test-kustomize-wasmtime:
	kustomize build --enable-alpha-plugins --enable-exec example-wasmtime
