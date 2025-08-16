@echo off
REM Generate Production Swagger Documentation
REM Usage: scripts\generate-production-swagger.bat [your-domain]

echo 🚀 Generating Production Swagger Documentation...

REM Check if swag is installed
swag version >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: swag command not found
    echo Please install swag first: go install github.com/swaggo/swag/cmd/swag@latest
    exit /b 1
)

REM Set production environment variables
set ENVIRONMENT=production
if "%1"=="" (
    set PROD_HOST=independent-carlene-tushar27x-a3461680.koyeb.app
) else (
    set PROD_HOST=%1
)

echo 📝 Using PROD_HOST: %PROD_HOST%

REM Generate Swagger docs
echo 🔧 Running swag init...
swag init -g cmd/main.go --parseDependency --parseInternal

echo ✅ Production Swagger documentation generated successfully!
echo 📚 Swagger UI will now use: %PROD_HOST%
echo.
echo 🔄 Next steps:
echo 1. Commit and push the updated docs/
echo 2. Redeploy to Koyeb
echo 3. Your Swagger UI will now send requests to %PROD_HOST%
