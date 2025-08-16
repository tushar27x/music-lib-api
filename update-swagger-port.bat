@echo off
REM Script to update Swagger port configuration
REM Usage: update-swagger-port.bat <port_number>

if "%1"=="" (
    echo Usage: %0 ^<port_number^>
    echo Example: %0 8082
    pause
    exit /b 1
)

set PORT=%1

echo Updating Swagger configuration to use port %PORT%...

REM Update main.go (using PowerShell for sed-like functionality)
powershell -Command "(Get-Content cmd/main.go) -replace 'localhost:[0-9]+', 'localhost:%PORT%' | Set-Content cmd/main.go"

REM Update README.md
powershell -Command "(Get-Content README.md) -replace 'localhost:[0-9]+', 'localhost:%PORT%' | Set-Content README.md"

REM Regenerate Swagger docs
echo Regenerating Swagger documentation...
swag init -g cmd/main.go

echo ✅ Swagger configuration updated to use port %PORT%
echo 🔄 Restart your server and access Swagger UI at: http://localhost:%PORT%/swagger/index.html
pause
