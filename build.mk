# Makefile

# Go source file
SRC = main.go

# Output directories
OUT_DIR = bin

# Executable names
WIN_EXEC = $(OUT_DIR)/main.exe
LINUX_EXEC = $(OUT_DIR)/main-linux
MAC_EXEC = $(OUT_DIR)/main-mac

# Build targets
all: windows linux mac

windows:
	GOOS=windows GOARCH=amd64 go build -o $(WIN_EXEC) $(SRC)

linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_EXEC) $(SRC)

mac:
	GOOS=darwin GOARCH=amd64 go build -o $(MAC_EXEC) $(SRC)

clean:
	rm -rf $(OUT_DIR)

.PHONY: all windows linux mac clean