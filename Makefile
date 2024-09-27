# Переменные для удобства
PROGRAM_NAME = myprogram
GO_FILES = ./client/cmd/main.go

# Цели для компиляции
all: build_linux build_macos build_windows

build_linux:
	GOOS=linux GOARCH=amd64 go build -o $(PROGRAM_NAME)_linux_amd64 $(GO_FILES)

build_macos:
	GOOS=darwin GOARCH=amd64 go build -o $(PROGRAM_NAME)_macos_amd64 $(GO_FILES)

build_windows:
	GOOS=windows GOARCH=amd64 go build -o $(PROGRAM_NAME)_windows_amd64.exe $(GO_FILES)

clean:
	rm -f $(PROGRAM_NAME)_linux_amd64 $(PROGRAM_NAME)_macos_amd64 $(PROGRAM_NAME)_windows_amd64.exe
