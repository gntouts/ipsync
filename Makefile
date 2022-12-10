.PHONY = default build docker go

define dockerstring
FROM golang:1.16-alpine

WORKDIR /app

COPY ./* ./
RUN go mod download && go mod tidy && go mod vendor

RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-w -extldflags "-static"' -o ipsync
endef
export dockerstring

SHELL := /bin/bash
DOCKER := $(shell command -v docker)
GOLANG := $(shell command -v go)

ifdef GOLANG
    build_target:=go
else
    build_target:=docker
endif

default: build

build: $(build_target)
	
docker:
	@echo "Building with Docker..."
ifdef DOCKER
		@echo "Creating Dockerfile"
		@echo "$$dockerstring" > Dockerfile.build
		@echo "Building Dockerfile"
		# @docker build -t ipsync:build -f Dockerfile.build .  > /dev/null 2>&1
		@docker build -t ipsync:build -f Dockerfile.build .
		@mkdir -p bin
		@echo "Building 'ipsync' binary"
		@docker run --name ipsync_build ipsync:build sh
		@docker cp ipsync_build:/app/ipsync ./bin/ipsync
		@docker rm -f ipsync_build > /dev/null 2>&1
		@rm -f Dockerfile.build > /dev/null 2>&1
		@docker rmi ipsync:build > /dev/null 2>&1
		@echo "Build complete."
		@echo "The compiled binary is placed in './bin/ipsync'"
else
		@echo "Docker is not installed, please install Go or Docker to build the binary."
endif


go:
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -o ./bin/ipsync
	@echo "Build complete."
	@echo "The compiled binary is placed in './bin/ipsync'"