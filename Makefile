version=$(shell git describe --tags --long --always|sed 's/^v//')

all: dep
	gom exec go build -ldflags "-X main.version=$(version)" uberstatus.go
	go fmt

dep:
	gom install
	# Hack around go's retarded way of dealing with "global" package naming
	mkdir -p _vendor/src/github.com/XANi
	ln -s . _vendor/src/github.com/XANi/uberstatus >/dev/null 2>&1 || true

gccgo: dep
	gom exec go build -compiler gccgo -gccgoflags "-O3" uberstatus.go
	go fmt
