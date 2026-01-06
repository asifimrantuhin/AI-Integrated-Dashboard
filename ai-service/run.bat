@echo off
setlocal enabledelayedexpansion

echo ========================================
echo   AI Service - Starting...
echo ========================================
echo.

REM Check if Python is available
python --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Python is not installed or not in PATH!
    echo Please install Python 3.8+ and try again.
    pause
    exit /b 1
)

REM Check if virtual environment exists
if not exist "venv\Scripts\python.exe" (
    echo Virtual environment not found. Creating one...
    python -m venv venv
    if errorlevel 1 (
        echo ERROR: Failed to create virtual environment!
        pause
        exit /b 1
    )
    echo Installing dependencies...
    call venv\Scripts\python.exe -m pip install --upgrade pip
    if errorlevel 1 (
        echo ERROR: Failed to upgrade pip!
        pause
        exit /b 1
    )
    call venv\Scripts\python.exe -m pip install -r requirements.txt
    if errorlevel 1 (
        echo ERROR: Failed to install dependencies!
        pause
        exit /b 1
    )
    echo Dependencies installed successfully!
    echo.
)

REM Check if all dependencies are installed (check for fastapi as indicator)
echo Checking if dependencies are installed...
venv\Scripts\python.exe -m pip show fastapi >nul 2>&1
if errorlevel 1 (
    echo Dependencies not found. Installing from requirements.txt...
    echo This may take a few minutes...
    call venv\Scripts\python.exe -m pip install --upgrade pip
    call venv\Scripts\python.exe -m pip install -r requirements.txt
    if errorlevel 1 (
        echo ERROR: Failed to install dependencies!
        pause
        exit /b 1
    )
    echo Dependencies installed successfully!
    echo.
)

REM Check if main.py exists
if not exist "main.py" (
    echo ERROR: main.py not found in current directory!
    echo Current directory: %CD%
    pause
    exit /b 1
)

REM Run the server using venv's Python directly
echo.
echo ========================================
echo   Starting FastAPI server...
echo ========================================
echo   Server URL: http://localhost:8888
echo   API Docs:   http://localhost:8888/docs
echo   Health:     http://localhost:8888/health
echo ========================================
echo.
echo Press Ctrl+C to stop the server
echo.

venv\Scripts\python.exe -m uvicorn main:app --reload --port 8888

if errorlevel 1 (
    echo.
    echo ERROR: Server failed to start!
    echo Check the error messages above.
    pause
    exit /b 1
)

