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
# run the server
@echo off

set BIN_PATH=C:\path\to\your\project\bin

echo Running the executable...
cd %BIN_PATH%
your_app.exe

pause
```


[1]: https://pkg.go.dev/badge/github.com/teocci/go-fiber-web.svg
[2]: https://pkg.go.dev/github.com/teocci/go-fiber-web
[3]: https://github.com/teocci/go-fiber-web/releases/tag/v1.0.0



