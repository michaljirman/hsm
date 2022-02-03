all: client server

protoc:
	@echo "Generating Go files"
	cd ./proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

server: protoc
	@echo "Building server"
	go build -o server \
		github.com/michaljirman/hsm/mpc

client: protoc
	@echo "Building client"
	go build -o client \
		github.com/michaljirman/hsm/sdk

clean:
	go clean github.com/michaljirman/hsm/...
	rm -f server client

.PHONY: client server protoc