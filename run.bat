@echo off
setlocal
set PROJECT_ROOT=%~dp0
set FRONTEND_DIR=%PROJECT_ROOT%frontend
set BACKEND_DIR=%PROJECT_ROOT%backend-api
set ADMIN_DIR=%PROJECT_ROOT%admin-cms
set AI_DIR=%PROJECT_ROOT%ai-service

echo ========================================
echo   Starting iDash Services
echo ========================================

echo Launching Backend API...
start "Backend API" cmd /k "cd /d "%BACKEND_DIR%" && go run main.go"

echo Launching Admin CMS (web)...
start "Admin CMS" cmd /k "cd /d "%ADMIN_DIR%" && php artisan serve --host=127.0.0.1 --port=8001"

echo Launching Admin CMS queue worker...
start "Admin CMS Queue" cmd /k "cd /d "%ADMIN_DIR%" && php artisan queue:work --tries=3"

echo Launching AI Service...
start "AI Service" cmd /k "cd /d "%AI_DIR%" && if not exist .venv (python -m venv .venv) && call .venv\Scripts\activate.bat && uvicorn main:app --host 0.0.0.0 --port 8000"

echo Launching Frontend (Vite)...
start "Frontend" cmd /k "cd /d "%FRONTEND_DIR%" && npm run dev -- --host"

echo Waiting for services to initialise...
timeout /t 8 >nul

start "" "http://localhost:5173"

echo All services started. Press Ctrl+C in each window to stop them.
endlocal
