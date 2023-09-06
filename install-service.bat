@echo off
REM Replace these values with your service-specific information
set SERVICE_NAME=wb-analyzer
set PROJECT_PATH=D:\code\go\%APP_NAME%
set BIN_PATH=%PROJECT_PATH%\bin
set SERVICE_PATH=%BIN_PATH%\%SERVICE_NAME%.exe
set SERVICE_DISPLAY_NAME=WB Analyzer
set SERVICE_DESCRIPTION=WB Analyzer service to run in background

REM Install the service
sc create "%SERVICE_NAME%" binPath= "%SERVICE_PATH%" DisplayName= "%SERVICE_DISPLAY_NAME%" description= "%SERVICE_DESCRIPTION%"
sc description "%SERVICE_NAME%" "%SERVICE_DESCRIPTION%"
sc start "%SERVICE_NAME%"

echo The service "%SERVICE_NAME%" has been installed and started.
echo To stop the service, use the following command:
echo sc stop "%SERVICE_NAME%"
echo To remove the service, use the following command:
echo sc delete "%SERVICE_NAME%"