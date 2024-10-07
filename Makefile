VERSION := $(shell git describe --tags --abbrev=0)

GORELEASER?=
export CGO_ENABLED?=0
PWD:=$(shell pwd)

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

test: contracts/node_modules
	@go test -coverprofile=coverage.txt -covermode=atomic -v ./...

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

## EGO and TEE bits
signed-build:
	docker run --rm -v $(PWD):/build -w /build -ti ghcr.io/edgelesssys/ego-dev /bin/bash -c "git config --global --add safe.directory /build && make ego-build"

# musl build:
# apt-get install -y --no-install-recommends musl-dev musl-tools
# CGO_ENABLED=1 CC=musl-gcc go build --ldflags '-linkmode=external -extldflags=-static' -o binary_name ./

ego-build: contracts/node_modules
	@ego-go build -v -ldflags '-linkmode=external -extldflags=-static' -ldflags "-X github.com/masa-finance/masa-oracle/internal/versioning.ApplicationVersion=${VERSION}" -o ./bin/masa-node ./cmd/masa-node
	@ego-go build -v -ldflags '-linkmode=external -extldflags=-static' -ldflags "-X github.com/masa-finance/masa-oracle/internal/versioning.ApplicationVersion=${VERSION}" -o ./bin/masa-node-cli ./cmd/masa-node-cli

sign:
	docker run --rm -v $(PWD):/build -w /build -ti ghcr.io/edgelesssys/ego-dev /bin/bash -c "ego sign ./tee/masa-node.json"

ego-run:
	docker run --rm -v $(PWD):/build -w /build -ti ghcr.io/edgelesssys/ego-dev /bin/bash -c "mkdir -p .masa && cp -rfv .env .masa && OE_SIMULATION=1 ego run ./bin/masa-node"