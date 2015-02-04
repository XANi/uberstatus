all: dep
	go fmt
	gom exec go build [a-z].*go

dep:
	gom install
