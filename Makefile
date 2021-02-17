TARGET=goskew

default: build

build:
	go build -v ./...

all: windows64 linux64 darwin64

windows64:
	env GOOS=windows GOARCH=amd64 go build -v -o windows-amd64/goskew.exe ./...
	zip -v goskew-windows-amd64.zip windows-amd64/goskew.exe LICENCE README.md

linux64:
	env GOOS=linux GOARCH=amd64 go build -v -o linux-amd64/goskew ./...
	tar -czvf goskew-linux-amd64.zip linux-amd64/goskew LICENCE README.md

darwin64:
	env GOOS=darwin GOARCH=amd64 go build -v -o darwin-amd64/goskew ./...
	env GOOS=darwin GOARCH=arm64 go build -v -o darwin-arm64/goskew ./...
	tar -czvf goskew-darwin-amd64.zip darwin-amd64/goskew LICENCE README.md
	tar -czvf goskew-darwin-arm64.zip darwin-arm64/goskew LICENCE README.md

test:
	go test -v ./...