build:
	@go build -v -o ./bin/masa-node ./cmd/masa-node
	@go build -v -o ./bin/masa-node-cli ./cmd/masa-node-cli

install:
	@sh ./node_install.sh
	
run: build
	@./bin/masa-node

client: build
	@./bin/masa-node-cli

test:
	@go test ./...

clean:
	@rm -rf bin
	@rm masa_oracle_node.log

wp:
	@pdflatex whitepaper.tex

proto:
	protoc --go_out=. --go_opt=paths=source_relative --proto_path=. pkg/proto/msg/message.proto

.PHONY: proto