.PHONY = default build go

SHELL := /bin/bash
GOLANG := $(shell command -v go)

ifdef GOLANG
    build_target:=go
else
    build_target:=no_go
endif

default: build

build: $(build_target)
	
no_go:
	@echo "Please install Go..."

go:
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -o ./bin/ipsync
	@echo "Build complete."
	@echo "The compiled binary is placed in './bin/ipsync'"