PROTO_DIR := proto

# Install necessary tools for protobuf compilation
install-tools:
	@echo "Installing necessary tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed successfully."

# Compile proto files
generate-proto: install-tools
	@echo "Generating gRPC code for UserService..."
	# Update PATH temporarily for the current shell session
	PATH=$(HOME)/go/bin:$(PATH) protoc -I=$(PROTO_DIR) \
		--go_out=. \
		--go-grpc_out=. \
		$(PROTO_DIR)/**/*.proto

download-dependencies:
	@echo "Downloading dependencies..."
	go mod tidy

# Run the entire pipeline
all: generate-proto download-dependencies
