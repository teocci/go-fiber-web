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
go build -o %BIN_PATH%\%APP_NAME%.exe

echo Copying files...
xcopy /s /i %VIEWS_PATH% %BIN_PATH%\views
xcopy /s /i %WEB_PATH% %BIN_PATH%\web

echo Build and copy completed.

pause