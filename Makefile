VERSION := $(shell git describe --tags --abbrev=0)

build:
	@go build -v -ldflags "-X github.com/masa-finance/masa-oracle/internal/constants.Version=$(VERSION)" -o ./bin/masa-node ./cmd/masa-node
	@go build -v -ldflags "-X github.com/masa-finance/masa-oracle/internal/constants.Version=$(VERSION)" -o ./bin/masa-node-cli ./cmd/masa-node-cli
	
install:
	@sh ./node_install.sh
	
run: build
	@./bin/masa-node

faucet: build
	./bin/masa-node --faucet

stake: build
	./bin/masa-node --stake 1000

client: build	
	@./bin/masa-node-cli

test:
	@go test ./...

clean:
	@rm -rf bin
	@rm masa_node.log
	@rm -rf ~/.masa/cache
	
proto:
	sh pkg/workers/messages/build.sh

docker-build:
	@docker build -t masa-node:latest .

docker-compose-up:
	@docker compose up --build

.PHONY: proto
