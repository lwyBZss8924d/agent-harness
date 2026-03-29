REPO_ROOT := $(CURDIR)
BUILD_DIR := $(REPO_ROOT)/build/bin

.PHONY: fmt build test ci clean install-dev install-release release-build release-install release-check

fmt:
	gofmt -w ./cmd ./internal

build:
	mkdir -p "$(BUILD_DIR)"
	go build -o "$(BUILD_DIR)/aih" ./cmd/aih
	go build -o "$(BUILD_DIR)/op-sa-broker" ./cmd/op-sa-broker
	go build -o "$(BUILD_DIR)/op-sa-broker-client" ./cmd/op-sa-broker-client

test:
	./scripts/test.sh

ci: test release-check

clean:
	rm -rf build dist

install-dev:
	./scripts/install-dev.sh

install-release: build
	./scripts/install-release.sh

release-build:
	./scripts/build-release.sh

release-install: release-build
	./scripts/install-release.sh

release-check:
	./scripts/check-release.sh
