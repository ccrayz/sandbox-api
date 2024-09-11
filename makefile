# Variables
PORT=8080
CALL_SERVICE=http://localhost:9999
BINARY_NAME=sandbox-api

# Run the application
run:
	PORT=$(PORT) \
	CALL_SERVICE=$(CALL_SERVICE) \
	go run main.go server

build:
	go build -o .build/$(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME)

.PHONY: run build clean
