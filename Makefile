build:
	@go build -v -o ./bin/masa-node ./cmd/masa-node

install:
	@sh ./node_install.sh
	
run: build
	@./bin/masa-node

test:
	@go test ./...

clean:
	@rm -rf bin
	@rm masa_oracle_node.log
