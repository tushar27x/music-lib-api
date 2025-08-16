@echo off
REM Script to generate Swagger documentation for different environments
REM Usage: generate-swagger.bat [environment]
REM Environments: development, staging, production

set ENVIRONMENT=%1
if "%ENVIRONMENT%"=="" set ENVIRONMENT=development

echo Generating Swagger documentation for environment: %ENVIRONMENT%

REM Set environment variables
set ENVIRONMENT=%ENVIRONMENT%

if "%ENVIRONMENT%"=="production" (
    set PROD_HOST=%PROD_HOST%
    if "%PROD_HOST%"=="" set PROD_HOST=api.yourdomain.com
    set STAGING_HOST=%STAGING_HOST%
    if "%STAGING_HOST%"=="" set STAGING_HOST=staging.yourdomain.com
    set DEV_HOST=%DEV_HOST%
    if "%DEV_HOST%"=="" set DEV_HOST=localhost:8082
    echo Production host: %PROD_HOST%
) else if "%ENVIRONMENT%"=="staging" (
    set PROD_HOST=%PROD_HOST%
    if "%PROD_HOST%"=="" set PROD_HOST=api.yourdomain.com
    set STAGING_HOST=%STAGING_HOST%
    if "%STAGING_HOST%"=="" set STAGING_HOST=staging.yourdomain.com
    set DEV_HOST=%DEV_HOST%
    if "%DEV_HOST%"=="" set DEV_HOST=localhost:8082
    if "%STAGING_HOST%"=="staging.yourdomain.com" (
        echo ⚠️  Warning: Using default staging host. Set STAGING_HOST env var for custom domain.
    )
    echo Staging host: %STAGING_HOST%
) else (
    set PROD_HOST=%PROD_HOST%
    if "%PROD_HOST%"=="" set PROD_HOST=api.yourdomain.com
    set STAGING_HOST=%STAGING_HOST%
    if "%STAGING_HOST%"=="" set STAGING_HOST=staging.yourdomain.com
    set DEV_HOST=%DEV_HOST%
    if "%DEV_HOST%"=="" set DEV_HOST=localhost:8082
    echo Development host: %DEV_HOST%
)

REM Generate Swagger docs
echo Running swag init...
swag init -g cmd/main.go

echo ✅ Swagger documentation generated for %ENVIRONMENT% environment
echo 📁 Documentation files created in docs/ directory
echo 🌐 Access Swagger UI at: http://%DEV_HOST%/swagger/index.html (dev) or https://%PROD_HOST%/swagger/index.html (prod)
pause
