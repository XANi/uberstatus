version=$(shell git describe --tags --long --always|sed 's/^v//')
binfile=uberstatus

all:
	go build -ldflags "-X main.version=$(version)" $(binfile).go
	-@go fmt

static: glide.lock vendor
	go build -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o $(binfile).static $(binfile).go

version:
	@echo $(version)
