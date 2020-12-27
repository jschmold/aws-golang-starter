.PHONY: build deploy clean e2e clean_mocks build_mocks verify_mocks mocks test migrator migrate_up migrate_down build_ci
.SILENT:

SHELL := bash

modules  := $(shell cd modules && ls)
services := $(shell cd services && ls)

indent := sed 's/^/    /';
build_cmd := env GOOS=linux go build -ldflags="-s -w"

build:
	for target in $(services) ; do \
		echo "Building $$target: " ;\
		cd services/$$target && (make | ${indent}) ;\
	done

build_ci:
	for target in $(services) ; do \
		echo "Building $$target: " ;\
		cd services/$$target ;\
		if ! make; then \
			printf "    Failed to build.\n" ;\
			exit 1 ;\
		fi; \
	done

clean:
	for target in $(services) ; do \
		printf "Clean $$target... "; \
		cd services/$$target && (make clean | ${indent}); \
		echo "Done!" ; \
	done

deploy: clean build
	for target in $(services) ; do \
		echo "Deploying $$target: "; \
		cd services/$$target && (sls deploy | ${indent}); \
	done

deploy_ci: clean build_ci
	for target in $(services) ; do \
		echo "Deploying $$target: "; \
		cd services/$$target && sls deploy; \
	done


test:
	go test ./modules/...

build_mocks:
	./scripts/build-mocks.sh

clean_mocks:
	printf "Cleaning mocks"
	rm -rf mocks
	echo " - Done!"

mocks: clean_mocks build_mocks
	echo "Done!"

migrator:
	cd migrations && make

migrate_up: migrator
	./bin/migrator up

migrate_down: migrator
	./bin/migrator down

e2e:
	go test ./test/...

