@echo off
SETLOCAL EnableDelayedExpansion

echo ==================================================
echo  Network-EDR Podman Quick Deployment (CMD)
echo ==================================================

REM Check if Podman is installed
where podman >nul 2>&1
if %errorlevel% neq 0 (
    echo [Error] Podman is not installed or not in PATH.
    exit /b 1
)

REM Check Podman machine status
for /f "tokens=*" %%i in ('podman machine inspect --format "{{.State}}" 2^>nul') do set MACHINE_STATE=%%i
if defined MACHINE_STATE (
    if not "%MACHINE_STATE%"=="running" (
        echo Starting Podman machine...
        podman machine start
    )
)

set COMPOSE_FILE=deploy\compose.prod.yml
if not exist %COMPOSE_FILE% (
    if exist compose.prod.yml (
        set COMPOSE_FILE=compose.prod.yml
    ) else (
        echo [Error] Compose file not found.
        exit /b 1
    )
)

echo [1/2] Building and starting containers with Podman...
where podman-compose >nul 2>&1
if %errorlevel% equ 0 (
    podman-compose -f %COMPOSE_FILE% up -d --build
) else (
    podman compose -f %COMPOSE_FILE% up -d --build
)

echo [2/2] Deployment complete!
echo Access the Network-EDR dashboard at http://localhost:8080
pause
