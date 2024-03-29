.PHONY: build

PROJECT= Mock Net

BINARY=mock-net

default:
	@echo ${PROJECT}

clean:
	@go clean
	@rm -rf tmp

check:
	@go fmt
	@go vet

build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./tmp/${BINARY}.exe main.go

macos-intel:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./tmp/${BINARY} main.go