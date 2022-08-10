version=$(shell git describe --tags --long --always|sed 's/^v//')
binfile=uberstatus

all:
	go build -ldflags "-X main.version=$(version)" $(binfile).go
	-@go fmt

arch:
	mkdir -p bin
	CGO_ENABLED=0 GOARCH=arm64 go build  -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o bin/$(binfile).amd64 $(binfile).go

static: glide.lock vendor
	go build -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o $(binfile).static $(binfile).go

version:
	@echo $(version)
