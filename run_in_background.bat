@echo off

set APP_NAME=wb-analyzer
set GO_SERVER_PATH=%APP_NAME%.exe

echo Starting Go Application...

REM Replace with the path to your application executable
start /b .\%GO_SERVER_PATH%

echo Go Application is running in the background.

pause