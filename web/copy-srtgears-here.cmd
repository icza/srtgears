rmdir github.com /s /q

mkdir github.com\icza\srtgears\exec
copy %GOPATH%\src\github.com\icza\srtgears\*.go github.com\icza\srtgears\
copy %GOPATH%\src\github.com\icza\srtgears\exec\*.go github.com\icza\srtgears\exec\
del github.com\icza\srtgears\*_test.go

pause