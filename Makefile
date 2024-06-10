build:
	@go mod tidy
	@go build -v -o ./bin/masa-node ./cmd/masa-node
	@go build -v -o ./bin/masa-node-cli ./cmd/masa-node-cli
	@go build -v -o ./bin/masa-cli ./cmd/masa-cli

install:
	@sh ./node_install.sh
	
run: build
	@./bin/masa-node

faucet: build
	./bin/masa-node --faucet

stake: build
	./bin/masa-node --stake 1000

client: build	
	@./bin/masa-cli

test:
	@go test ./...

clean:
	@rm -rf bin
	@rm masa_node.log

wp:
	@pdflatex whitepaper.tex

proto:
	sh pkg/workers/messages/build.sh

.PHONY: proto

