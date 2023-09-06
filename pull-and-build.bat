@echo off

set APP_NAME=wb-analyzer
set PROJECT_PATH=D:\code\go\%APP_NAME%
set BIN_PATH=%PROJECT_PATH%\bin
set VIEWS_PATH=%PROJECT_PATH%\views
set WEB_PATH=%PROJECT_PATH%\web

set BIN_VIEWS_PATH=%BIN_PATH%\views
set BIN_WEB_PATH=%BIN_PATH%\web

echo Pulling from Git...
cd %PROJECT_PATH%
git pull

echo Building Go project...
go build -o %APP_NAME%.exe

if %errorlevel% neq 0 (
    echo Error: Build failed
    pause
    exit /b %errorlevel%
)

if not exist %BIN_PATH% (
    mkdir %BIN_PATH%
)

echo Copying files...
xcopy /y /d /i %APP_NAME%.exe %BIN_PATH%
xcopy /y /d /i .env %BIN_PATH%
xcopy /y /d /i config.json %BIN_PATH%
xcopy /y /d /i run_in_background.bat %BIN_PATH%
xcopy /y /d /i %VIEWS_PATH%\* %BIN_VIEWS_PATH%
xcopy /y /d /i %WEB_PATH%\* %BIN_WEB_PATH%

echo Build and copy completed.

pause