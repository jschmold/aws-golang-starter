.PHONY: build clean deploy
SHELL := bash

funcs := register
build_cmd := env GOOS=linux go build -ldflags="-s -w"

build:
	for target in $(funcs); do \
		printf "Service ./$$target/main.go %10s";\
		if $(build_cmd) -o bin/$$target ./$$target/main.go; \
		then \
			printf "Status: Success" ;\
		else \
			printf "Status: Failed" ;\
			printf "\n" ;\
			exit 1;\
			break ;\
		fi ;\
		printf "\n" ;\
	done

clean:
	rm -rf ./bin
