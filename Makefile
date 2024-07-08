windows:
	@GOOS=windows go build -o bin/cpkgs_windows.exe main.go
linux:
	@GOOS=linux go build -o bin/cpkgs_linux main.go
darwin:
	@GOOS=darwin go build -o bin/cpkgs_darwin main.go
all: windows linux darwin