@echo off
REM Check if Go code is properly formatted

echo Checking Go code formatting...
gofmt -l . > unformatted.txt

for /f %%A in ("unformatted.txt") do set size=%%~zA

if %size% GTR 0 (
    echo The following files are not properly formatted:
    type unformatted.txt
    echo Run 'gofmt -w .' to format all files
    del unformatted.txt
    exit /b 1
) else (
    echo All Go files are properly formatted
    del unformatted.txt
    exit /b 0
)
