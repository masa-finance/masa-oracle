VERSION := $(shell git describe --tags --abbrev=0)

GORELEASER?=
CGO_ENABLED?=0

# check if goreleaser exists
ifeq (, $(shell which goreleaser))
	GORELEASER=curl -sfL https://goreleaser.com/static/run | bash -s --
else
	GORELEASER=$(shell which goreleaser)
endif

print-version:
	@echo "Version: ${VERSION}"

contracts/node_modules:
	@go generate ./...

dev-dist:
	$(GORELEASER) build --snapshot --single-target --clean

dist:
	$(GORELEASER) build --single-target --clean

build: contracts/node_modules
	@go build -v -ldflags "-X github.com/masa-finance/masa-oracle/internal/versioning.ApplicationVersion=${VERSION}" -o ./bin/masa-node ./cmd/masa-node
	@go build -v -ldflags "-X github.com/masa-finance/masa-oracle/internal/versioning.ApplicationVersion=${VERSION}" -o ./bin/masa-node-cli ./cmd/masa-node-cli

install:
	@sh ./node_install.sh

run: build
	@./bin/masa-node

run-api-enabled: build
	@./bin/masa-node --api-enabled=true

faucet: build
	./bin/masa-node --faucet

stake: build
	./bin/masa-node --stake 1000

client: build	
	@./bin/masa-node-cli

# TODO: Add -race and fix race conditions
test: contracts/node_modules
	@go test -coverprofile=coverage.txt -covermode=atomic -v -count=1 -shuffle=on ./...

ci-lint:
	go mod tidy && git diff --exit-code
	go mod download
	go mod verify
	gofmt -s -w . && git diff --exit-code
	go vet ./...
	golangci-lint run

clean:
	@rm -rf bin

	@if [ -d ~/.masa/blocks ]; then rm -rf ~/.masa/blocks; fi
	@if [ -d ~/.masa/cache ]; then rm -rf ~/.masa/cache; fi	
	@if [ -f masa_node.log ]; then rm masa_node.log; fi

proto:
	sh pkg/workers/messages/build.sh

docker-build:
	@docker build -t masa-node:latest .

docker-compose-up:
	@docker compose up --build

.PHONY: proto

## Test e2e stack
e2e-start:
	bash tests/e2e/start.sh

e2e-prepare:
	bash tests/e2e/prepare.sh

e2e-clean:
	bash tests/e2e/clean.sh

e2e-logs:
	cd tests/e2e/work && docker-compose logs -f

e2e-run:
	E2E=true go test -v ./tests/e2e

e2e-restart: e2e-clean e2e-prepare e2e-start e2e-run

e2e-reprepare: e2e-clean e2e-prepare e2e-start e2e-run
