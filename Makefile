windows:
	@GOOS=windows go build -o bin/cpkgs_windows.exe src/main.go
linux:
	@GOOS=linux go build -o bin/cpkgs_linux src/main.go
darwin:
	@GOOS=darwin go build -o bin/cpkgs_darwin src/main.go
all: windows linux darwin