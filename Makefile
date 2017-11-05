package = github.com/jbarone/ghostToHugo

.PHONY: release

release: package
	rm -rf dist

package: build
	mkdir -p release
	tar -cvzf release/ghostToHugo-Darwin-x86_64.tar.gz --directory=dist/darwin/amd64 .
	tar -cvzf release/ghostToHugo-Linux-x86_64.tar.gz --directory=dist/linux/amd64 .
	zip -j release/ghostToHugo-Windows-x86_64 dist/windows/amd64/*
	tar -cvzf release/ghostToHugo-Linux-i386.tar.gz --directory=dist/linux/386 .
	zip -j release/ghostToHugo-Windows-i386 dist/windows/386/*

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/darwin/amd64/ghostToHugo $(package)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/linux/amd64/ghostToHugo $(package)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/windows/amd64/ghostToHugo.exe $(package)
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o dist/linux/386/ghostToHugo $(package)
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o dist/windows/386/ghostToHugo.exe $(package)
