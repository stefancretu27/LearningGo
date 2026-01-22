BUILD_DIR := ./build
EXE_NAME := learningGo

.PHONY: all clean build deploy run stop

all: build

build:
	@echo "Building app: $@"
	mkdir build
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(EXE_NAME) ./*.go

clean:
	@echo "Cleaning build directory"
	rm -rf $(BUILD_DIR)