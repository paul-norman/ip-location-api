BINARY_NAME=ip-location-api

ifeq ($(OS), Windows_NT) 
	DETECTED_OS := Windows
	BUILD_COMMAND := build_windows
	CLEAN_COMMAND := clean_windows
else
	DETECTED_OS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
	BUILD_COMMAND := build_other
	CLEAN_COMMAND := clean_other
endif
 
all: test build

# Build commands
build: update $(BUILD_COMMAND)

build_linux_amd64:
	GOARCH=amd64 GOOS=linux go build -o builds/$(BINARY_NAME)-linux-x64.bin .

build_windows_amd64:
	GOARCH=amd64 GOOS=windows go build -o builds/$(BINARY_NAME)-windows-x64.exe .

build_darwin_amd64:
	GOARCH=amd64 GOOS=darwin go build -o builds/$(BINARY_NAME)-darwin-x64.dmg .

build_linux_arm64:
	GOARCH=arm64 GOOS=linux go build -o builds/$(BINARY_NAME)-linux-arm64.bin .

build_windows_arm64:
	GOARCH=arm64 GOOS=windows go build -o builds/$(BINARY_NAME)-windows-arm64.exe .

build_darwin_arm64:
	GOARCH=arm64 GOOS=darwin go build -o builds/$(BINARY_NAME)-darwin-arm64.dmg .

build_other: build_linux_amd64 build_windows_amd64 build_darwin_amd64
build_other: build_linux_arm64 build_windows_arm64 build_darwin_arm64

build_windows:
	set "GOARCH=amd64" && set "GOOS=linux" && go build -o builds\$(BINARY_NAME)-linux-x64.bin .
	set "GOARCH=amd64" && set "GOOS=windows" && go build -o builds\$(BINARY_NAME)-windows-x64.exe .
	set "GOARCH=amd64" && set "GOOS=darwin" && go build -o builds\$(BINARY_NAME)-darwin-x64.dmg .
	set "GOARCH=arm64" && set "GOOS=linux" && go build -o builds\$(BINARY_NAME)-linux-arm64.bin .
	set "GOARCH=arm64" && set "GOOS=windows" && go build -o builds\$(BINARY_NAME)-windows-arm64.exe .
	set "GOARCH=arm64" && set "GOOS=darwin" && go build -o builds\$(BINARY_NAME)-darwin-arm64.dmg .

# Clean commands
clean: $(CLEAN_COMMAND)

clean_go:
	go clean

clean_other: clean_go
	rm builds/$(BINARY_NAME)-linux-x64.bin
	rm builds/$(BINARY_NAME)-windows-x64.exe
	rm builds/$(BINARY_NAME)-darwin-x64.dmg
	rm builds/$(BINARY_NAME)-linux-arm64.bin
	rm builds/$(BINARY_NAME)-windows-arm64.exe
	rm builds/$(BINARY_NAME)-darwin-arm64.dmg

clean_windows: clean_go
	del "builds\$(BINARY_NAME)-linux-x64.bin"
	del "builds\$(BINARY_NAME)-windows-x64.exe"
	del "builds\$(BINARY_NAME)-darwin-x64.dmg"
	del "builds\$(BINARY_NAME)-linux-arm64.bin"
	del "builds\$(BINARY_NAME)-windows-arm64.exe"
	del "builds\$(BINARY_NAME)-darwin-arm64.dmg"

# Dev commands
dev: update
	go run .

# Run commands
run: clean build
	./builds/$(BINARY_NAME)

# Test commands
test: update
	go test .

# Update commands
update:
	go get -u
	go mod tidy