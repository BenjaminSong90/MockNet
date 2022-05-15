.PHONY: build

PROJECT= Mock Net

BINARY=MockNet

default:
	@echo ${PROJECT}

clean:
	@go clean
	rm ${BINARY}*

check:
	@go fmt
	@go vet

build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}.exe main.go

build-macos-intel:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BINARY} main.go