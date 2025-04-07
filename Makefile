.PHONY: build run-server run-client test test-sound clean all

all: build

build:
	@echo "Building ECG Tool..."
	@mkdir -p bin
	go build -o bin/ecg-server cmd/server/main.go
	go build -o bin/ecg-client cmd/client/main.go
	go build -o bin/sound-test cmd/sound_test/main.go
	@echo "Build complete. Binaries are in the bin/ directory"

run-server:
	@echo "Starting ECG server..."
	go run cmd/server/main.go

run-client:
	@echo "Starting ECG client..."
	go run cmd/client/main.go

test:
	@echo "Running tests..."
	go test ./...

test-sound:
	@echo "Testing sound functionality..."
	go run cmd/sound_test/main.go

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf logs/
