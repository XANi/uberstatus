all: dep
	gom exec go build [a-z]*go
	go fmt

dep:
	gom install
	# Hack around go's retarded way of dealing with "global" package naming
	mkdir -p _vendor/src/github.com/XANi
	ln -s . _vendor/src/github.com/XANi/uberstatus >/dev/null 2>&1 || true
