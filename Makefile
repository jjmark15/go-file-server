BINARY_NAME=goserver
INSTALL_PATH=/usr/local/bin

.PHONY: all build clean install uninstall

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)

uninstall:
	@echo "Removing from $(INSTALL_PATH)..."
	sudo rm $(INSTALL_PATH)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

run:
	@go run main.go
