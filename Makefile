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
	@mkdir -p dist/bin
	@go mod tidy
	@CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -o ./dist/bin/ipsync
	@echo "Build complete. The compiled binary is placed in './dist/bin/ipsync'"

image: clean
	@go mod tidy
	@mkdir -p dist/src
	@cp main.go dist/src/
	@cp utils.go dist/src
	@cp go.* dist/src
	@cp -r pkg dist/src
	docker build -t gntouts/ipsync:latest -f dist/Dockerfile dist
	@rm -fr dist/src

clean:
	@rm -fr dist/src
	@rm -fr dist/bin

compose:
	cd dist && docker compose up -d