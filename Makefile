# This file is generated with create-go-app: do not edit.
.PHONY: build clean test snapshot all help

# special target to export all variables
.EXPORT_ALL_VARIABLES:

## build: build the binaries only
build:
	goreleaser build --rm-dist --snapshot

## snapshot: create a snapshot release
snapshot:
	goreleaser release --snapshot --rm-dist --skip-sign

## clean: cleanup
clean:
	rm -rf ./dist

all: build

## test: run linter and tests
test:
	go generate ./...
	golangci-lint run
	go test -v -count=1 ./...

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
