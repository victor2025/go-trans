# README for build
## build sample
```makefile
OUTPUT_DIR=build
GO_FILE=example.go

all: linux_amd64 linux_arm64 linux_arm windows_amd64 macos_arm64

linux_amd64:
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o $(OUTPUT_DIR)/libexample_linux_amd64.so $(GO_FILE)

linux_arm64:
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -buildmode=c-shared -o $(OUTPUT_DIR)/libexample_linux_arm64.so $(GO_FILE)

linux_arm:
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -buildmode=c-shared -o $(OUTPUT_DIR)/libexample_linux_arm.so $(GO_FILE)

windows_amd64:
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o $(OUTPUT_DIR)/example_windows_amd64.dll $(GO_FILE)

macos_arm64:
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=arm64 go build -buildmode=c-shared -o $(OUTPUT_DIR)/libexample_darwin_arm64.dylib $(GO_FILE)

clean:
	rm -rf $(OUTPUT_DIR)
```