set VERSION=1.1

@set GOBUILD=go build -ldflags "-X main.Version=%VERSION%" -o

@set GOOS=windows
@set GOARCH=amd64
%GOBUILD% srtgears-%VERSION%-%GOOS%-%GOARCH%/srtgears.exe

@set GOARCH=386
%GOBUILD% srtgears-%VERSION%-%GOOS%-%GOARCH%/srtgears.exe

@set GOOS=linux
%GOBUILD% srtgears-%VERSION%-%GOOS%-%GOARCH%/srtgears

@set GOARCH=amd64
%GOBUILD% srtgears-%VERSION%-%GOOS%-%GOARCH%/srtgears

@set GOOS=darwin
%GOBUILD% srtgears-%VERSION%-%GOOS%-%GOARCH%/srtgears

@set GOARCH=386
%GOBUILD% srtgears-%VERSION%-%GOOS%-%GOARCH%/srtgears

@set GOOS=windows
@set GOARCH=amd64

@echo:
@echo DONE! Press a key.
@pause>nul
