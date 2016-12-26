version=$(shell git describe --tags --long --always|sed 's/^v//')
binfile=uberstatus

all: glide.lock vendor
	rm -rf _vendor
	go build -ldflags "-X main.version=$(version)" $(binfile).go
	-@go fmt

static: glide.lock vendor
	go build -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o $(binfile).static $(binfile).go

clean:
	rm -rf vendor
	rm -rf _vendor
vendor: glide.lock
	glide install && touch vendor
glide.lock: glide.yaml
	glide update && touch glide.lock
glide.yaml:
version:
	@echo $(version)
