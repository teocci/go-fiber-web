@echo off

set "project_dir=D:\Apps\wb-analyzer"  REM Replace with your project directory
set "bin_dir=%project_dir%\bin"
set "go_exe=go.exe"
set "git_exe=git.exe"

echo Pulling from Git...
cd /d "%project_dir%"
%git_exe% pull origin master

echo Building the Go project...
%go_exe% build -o %bin_dir%\wb-analyzer.exe

echo Copying the executable to the bin directory...
copy /y %bin_dir%\wb-analyzer.exe %bin_dir%\

echo Build process completed.
pause