if NOT EXIST "app.yaml" goto error
if NOT EXIST "sgw/sgw.go" goto error

rmdir github.com /s /q

mkdir github.com\icza\srtgears\exec
copy %GOPATH%\src\github.com\icza\srtgears\*.go github.com\icza\srtgears\
copy %GOPATH%\src\github.com\icza\srtgears\exec\*.go github.com\icza\srtgears\exec\
del github.com\icza\srtgears\*_test.go

goto end

:error
echo It looks the script was launched from improper folder. Aborting...

:end
pause
