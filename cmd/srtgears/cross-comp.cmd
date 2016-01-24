set version=1.0

set GOOS=windows
set GOARCH=amd64
go build -o srtgears-%version%-%GOOS%-%GOARCH%/srtgears.exe

set GOARCH=386
go build -o srtgears-%version%-%GOOS%-%GOARCH%/srtgears.exe

set GOOS=linux
go build -o srtgears-%version%-%GOOS%-%GOARCH%/srtgears

set GOARCH=amd64
go build -o srtgears-%version%-%GOOS%-%GOARCH%/srtgears

set GOOS=darwin
go build -o srtgears-%version%-%GOOS%-%GOARCH%/srtgears

set GOARCH=386
go build -o srtgears-%version%-%GOOS%-%GOARCH%/srtgears

set GOOS=windows
set GOARCH=amd64
