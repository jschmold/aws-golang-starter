.PHONY: build
.SILENT:

SHELL := bash

build_cmd := env GOOS=linux go build -ldflags="-s -w"

build: main.go cmd.go
	${build_cmd} -o ../bin/migrator $+
