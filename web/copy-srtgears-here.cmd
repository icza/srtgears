rmdir github.com /s /q

mkdir github.com\gophergala2016\srtgears\exec
copy %GOPATH%\src\github.com\gophergala2016\srtgears\*.go github.com\gophergala2016\srtgears\
copy %GOPATH%\src\github.com\gophergala2016\srtgears\exec\*.go github.com\gophergala2016\srtgears\exec\
del github.com\gophergala2016\srtgears\*_test.go

pause