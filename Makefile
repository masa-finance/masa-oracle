ifndef VERSION
$(error VERSION is not set. Please set it in your .env file)
endif

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
	@echo "node_modules already installed in Docker build"

swagger:
	@echo "Installing swag if needed..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Running swag init..."
	@swag init -g pkg/api/routes.go || echo "Swag init failed with status: $$?"

dev-dist:
	$(GORELEASER) build --snapshot --single-target --clean

dist:
	$(GORELEASER) build --single-target --clean

build: contracts/node_modules swagger
	@pwd
	@echo "Current directory contents:"
	@ls -la
	@echo "Contracts directory contents:"
	@ls -la contracts/
	@echo "Building masa-node..."
	@go build -v -ldflags "-X github.com/masa-finance/masa-oracle/internal/versioning.ApplicationVersion=${VERSION}" -o ./bin/masa-node ./cmd/masa-node
	@echo "Building masa-node-cli..."
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
