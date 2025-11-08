@echo off
setlocal
set PROJECT_ROOT=%~dp0
set FRONTEND_DIR=%PROJECT_ROOT%frontend
set BACKEND_DIR=%PROJECT_ROOT%backend-api
set ADMIN_DIR=%PROJECT_ROOT%admin-cms
set AI_DIR=%PROJECT_ROOT%ai-service

echo ========================================
echo   iDash Setup Wizard (Windows)
echo ========================================

REM --- Copy environment templates if missing ---
if not exist "%FRONTEND_DIR%\.env" (
    if exist "%FRONTEND_DIR%\.env.example" (
        copy "%FRONTEND_DIR%\.env.example" "%FRONTEND_DIR%\.env" >nul
        echo Created frontend .env from template.
    )
)

if not exist "%ADMIN_DIR%\.env" (
    if exist "%ADMIN_DIR%\.env.example" (
        copy "%ADMIN_DIR%\.env.example" "%ADMIN_DIR%\.env" >nul
        echo Created admin CMS .env from template.
    )
)

if not exist "%AI_DIR%\.env" (
    if exist "%AI_DIR%\.env.example" (
        copy "%AI_DIR%\.env.example" "%AI_DIR%\.env" >nul
        echo Created AI service .env from template.
    )
)

echo.
echo Installing frontend dependencies...
pushd "%FRONTEND_DIR%"
call npm install
popd

echo.
echo Installing backend API dependencies...
pushd "%BACKEND_DIR%"
call go mod download
popd

echo.
echo Installing admin CMS dependencies...
pushd "%ADMIN_DIR%"
if not exist "%ADMIN_DIR%\vendor" (
    call composer install
) else (
    call composer install --no-interaction
)
call php artisan key:generate --force
call php artisan migrate --force
popd

echo.
echo Preparing AI service virtual environment...
pushd "%AI_DIR%"
if not exist ".venv" (
    call python -m venv .venv
)
if exist ".venv\Scripts\python.exe" (
    call ".venv\Scripts\python.exe" -m pip install --upgrade pip
    call ".venv\Scripts\python.exe" -m pip install -r requirements.txt
) else (
    echo WARNING: Python virtual environment not found. Ensure Python 3.10+ is installed and rerun setup.
)
popd

echo.
echo ========================================
echo   Setup complete!
echo   Next steps:

echo   1. Verify database credentials inside admin-cms\.env
echo   2. Ensure MySQL is running before launching services.
echo   3. Use run.bat to start the stack.

echo ========================================
endlocal
pause
