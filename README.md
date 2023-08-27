## go-fiber-web [![Go Reference][1]][2]

`go-fiber-web` is an open-source sample using fiber as a webserver as basic clean project.

## Disclaimer
> This tool is limited to security research only, and the user assumes all legal and related responsibilities arising from its use! The author assumes no legal responsibility!

## config
Change the port for the web server on the `config.json` file
```json
{
  "web": {
    "port": 9012
  }
}
```
add a .env file with the following variables

```env
MPS_API_SECRET=<HASH>
```

## Run
```bash
go run main.go

open http://localhost:9012/page.html
```

### Scripts
```bash
@echo off

set APP_NAME=wb-analyzer
set PROJECT_PATH=D:\code\go\%APP_NAME%
set BIN_PATH=%PROJECT_PATH%\bin
set VIEWS_PATH=%PROJECT_PATH%\views
set WEB_PATH=%PROJECT_PATH%\web

echo Pulling from Git...
cd %PROJECT_PATH%
git pull

echo Building Go project...
go build -o %APP_NAME%.exe

echo Copying files...
xcopy /y /d /i %APP_NAME%.exe %BIN_PATH%
xcopy /y /d /i %VIEWS_PATH%\* %BIN_PATH%
xcopy /y /d /i %WEB_PATH%\* %BIN_PATH%

echo Build and copy completed.

pause
```


```bash
# run the server
@echo off

set APP_NAME=wb-analyzer
set BIN_PATH=D:\code\go\%APP_NAME%\bin

echo Running the executable...
cd %BIN_PATH%
%APP_NAME%.exe

pause
```


[1]: https://pkg.go.dev/badge/github.com/teocci/go-fiber-web.svg
[2]: https://pkg.go.dev/github.com/teocci/go-fiber-web
[3]: https://github.com/teocci/go-fiber-web/releases/tag/v1.0.0



