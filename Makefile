BINARY_NAME=ip-location-api

ifeq ($(OS), Windows_NT)
	DETECTED_OS := Windows
	BUILD_COMMAND := build_on_windows_all
	BUILD_COMMAND_LINUX_X64 := build_on_windows_linux_x64
	BUILD_COMMAND_LINUX_ARM64 := build_on_windows_linux_arm64
	BUILD_COMMAND_DARWIN_X64 := build_on_windows_darwin_x64
	BUILD_COMMAND_DARWIN_ARM64 := build_on_windows_darwin_arm64
	BUILD_COMMAND_WINDOWS_X64 := build_on_windows_windows_x64
	BUILD_COMMAND_WINDOWS_ARM64 := build_on_windows_windows_arm64
	CLEAN_COMMAND := clean_on_windows
else
	DETECTED_OS := $(shell sh -c 'uname 2>/dev/null || echo Unknown')
	BUILD_COMMAND := build_on_unix_all
	BUILD_COMMAND_LINUX_X64 := build_on_unix_linux_x64
	BUILD_COMMAND_LINUX_ARM64 := build_on_unix_linux_arm64
	BUILD_COMMAND_DARWIN_X64 := build_on_unix_darwin_x64
	BUILD_COMMAND_DARWIN_ARM64 := build_on_unix_darwin_arm64
	BUILD_COMMAND_WINDOWS_X64 := build_on_unix_windows_x64
	BUILD_COMMAND_WINDOWS_ARM64 := build_on_unix_windows_arm64
	CLEAN_COMMAND := clean_on_unix
endif

all: test build

# Build commands
build: update $(BUILD_COMMAND)
build_linux: update $(BUILD_COMMAND_LINUX_X64) $(BUILD_COMMAND_LINUX_ARM64)
build_linux_x64: update $(BUILD_COMMAND_LINUX_X64)
build_linux_arm64: update $(BUILD_COMMAND_LINUX_ARM64)
build_darwin: update $(BUILD_COMMAND_DARWIN_X64) $(BUILD_COMMAND_DARWIN_ARM64)
build_darwin_x64: update $(BUILD_COMMAND_DARWIN_X64)
build_darwin_arm64: update $(BUILD_COMMAND_DARWIN_ARM64)
build_windows: update $(BUILD_COMMAND_WINDOWS_X64) $(BUILD_COMMAND_WINDOWS_ARM64)
build_windows_x64: update $(BUILD_COMMAND_WINDOWS_X64)
build_windows_arm64: update $(BUILD_COMMAND_WINDOWS_ARM64)

# Build on Unix
build_on_unix_all: build_on_unix_linux_x64 build_on_unix_linux_arm64 build_on_unix_darwin_x64 build_on_unix_darwin_arm64 build_on_unix_windows_x64 build_on_unix_windows_arm64
build_on_unix_linux_x64:
	GOARCH=amd64 GOOS=linux go build -o builds/$(BINARY_NAME)-linux-x64.bin .
build_on_unix_linux_arm64:
	GOARCH=arm64 GOOS=linux go build -o builds/$(BINARY_NAME)-linux-arm64.bin .
build_on_unix_darwin_x64:
	GOARCH=amd64 GOOS=darwin go build -o builds/$(BINARY_NAME)-darwin-x64.dmg .
build_on_unix_darwin_arm64:
	GOARCH=arm64 GOOS=darwin go build -o builds/$(BINARY_NAME)-darwin-arm64.dmg .
build_on_unix_windows_x64:
	GOARCH=amd64 GOOS=windows go build -o builds/$(BINARY_NAME)-windows-x64.exe .
build_on_unix_windows_arm64:
	GOARCH=arm64 GOOS=windows go build -o builds/$(BINARY_NAME)-windows-arm64.exe .

# Build on Windows
build_on_windows_all: build_on_windows_linux_x64 build_on_windows_linux_arm64 build_on_windows_darwin_x64 build_on_windows_darwin_arm64 build_on_windows_windows_x64 build_on_windows_windows_arm64
build_on_windows_linux_x64:
	set "GOARCH=amd64" && set "GOOS=linux" && go build -o builds\$(BINARY_NAME)-linux-x64.bin .
build_on_windows_linux_arm64:
	set "GOARCH=arm64" && set "GOOS=linux" && go build -o builds\$(BINARY_NAME)-linux-arm64.bin .
build_on_windows_darwin_x64:
	set "GOARCH=amd64" && set "GOOS=darwin" && go build -o builds\$(BINARY_NAME)-darwin-x64.dmg .
build_on_windows_darwin_arm64:
	set "GOARCH=arm64" && set "GOOS=darwin" && go build -o builds\$(BINARY_NAME)-darwin-arm64.dmg .
build_on_windows_windows_x64:
	set "GOARCH=amd64" && set "GOOS=windows" && go build -o builds\$(BINARY_NAME)-windows-x64.exe .
build_on_windows_windows_arm64:
	set "GOARCH=arm64" && set "GOOS=windows" && go build -o builds\$(BINARY_NAME)-windows-arm64.exe .

# Clean commands
clean: $(CLEAN_COMMAND)

clean_go:
	go clean

clean_on_unix: clean_go
	rm -f builds/$(BINARY_NAME)-linux-x64.bin
	rm -f builds/$(BINARY_NAME)-windows-x64.exe
	rm -f builds/$(BINARY_NAME)-darwin-x64.dmg
	rm -f builds/$(BINARY_NAME)-linux-arm64.bin
	rm -f builds/$(BINARY_NAME)-windows-arm64.exe
	rm -f builds/$(BINARY_NAME)-darwin-arm64.dmg

clean_on_windows: clean_go
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

# Docker commands
docker:
	docker build -t ip-location-api .