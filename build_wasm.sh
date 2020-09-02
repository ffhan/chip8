#!/bin/bash

GOOS=js GOARCH=wasm go build -o examples/wasm/main.wasm examples/wasm/wasm.go
