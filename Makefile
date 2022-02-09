all: mpc-coordinator mpc-signer

protoc:
	@echo "Generating Go files"
	cd ./proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

mpc-signer: protoc
	@echo "Building MPC Signer"
	ego-go build -o signer \
	    github.com/michaljirman/hsm/mpc-signer && ego sign mpc-signer

mpc-signer-dev: protoc
	@echo "Building MPC Signer for localhost development"
	go build -o signer github.com/michaljirman/hsm/mpc-signer

mpc-coordinator: protoc
	@echo "Building MPC Coordinator"
	go build -o coordinator \
		github.com/michaljirman/hsm/mpc-coordinator

clean:
	go clean github.com/michaljirman/hsm/...
	rm -f signer coordinator

.PHONY: mpc-coordinator mpc-signer mpc-signer-dev protoc