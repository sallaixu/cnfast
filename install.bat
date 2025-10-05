@echo off
REM cnfast Windows Install Script
REM For Windows 10/11

chcp 65001 >nul 2>&1
setlocal enabledelayedexpansion

echo ================================================
echo      cnfast Windows Install Script
echo ================================================
echo.

REM Define variables
set "BASE_URL=https://gitee.com/sallai/cnfast/releases/download/latest"
set "BINARY_NAME=cnfast.exe"
set "INSTALL_DIR=%LOCALAPPDATA%\cnfast"

REM Detect system architecture
echo [1/5] Detecting system architecture...
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    echo       Detected: AMD64
    set "ARCH_SUFFIX=amd64"
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    echo       Detected: ARM64
    set "ARCH_SUFFIX=arm64"
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    echo       Detected: x86
    set "ARCH_SUFFIX=386"
) else (
    echo       Error: Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

REM Build download URL
set "DOWNLOAD_URL=%BASE_URL%/cnfast-windows-%ARCH_SUFFIX%.exe"
echo       URL: %DOWNLOAD_URL%
echo.

REM Create temp file
set "TMP_FILE=%TEMP%\cnfast_%RANDOM%.exe"

REM Download file
echo [2/5] Downloading cnfast...
echo       Please wait, this may take a few minutes...
echo.

REM Use PowerShell to download
powershell -NoProfile -ExecutionPolicy Bypass -Command "[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; $ProgressPreference = 'SilentlyContinue'; try { Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%TMP_FILE%' -UseBasicParsing; Write-Host '       Download completed!' -ForegroundColor Green } catch { Write-Host '       Download failed' -ForegroundColor Red; exit 1 }"

if errorlevel 1 (
    echo       Error: Download failed, please check your network connection
    exit /b 1
)

REM Check file exists
if not exist "%TMP_FILE%" (
    echo       Error: Download failed, file does not exist
    exit /b 1
)

echo.
echo [3/5] Installing to %INSTALL_DIR%...

REM Create install directory
if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)

REM Delete old version if exists
if exist "%INSTALL_DIR%\%BINARY_NAME%" (
    del /f /q "%INSTALL_DIR%\%BINARY_NAME%" 2>nul
    if errorlevel 1 (
        echo       Warning: Cannot delete old version, may be in use
        echo       Please close all cnfast processes and try again
        del /f /q "%TMP_FILE%" 2>nul
        exit /b 1
    )
)

REM Move file to install directory
move /y "%TMP_FILE%" "%INSTALL_DIR%\%BINARY_NAME%" >nul

if errorlevel 1 (
    echo       Error: Installation failed
    exit /b 1
)

echo       Installation completed!
echo.

REM Add to PATH
echo [4/5] Configuring environment variables...
echo.

REM Use PowerShell to add to user PATH
powershell -NoProfile -ExecutionPolicy Bypass -Command "$userPath = [Environment]::GetEnvironmentVariable('Path', 'User'); if ($userPath -notlike '*%INSTALL_DIR%*') { $newPath = $userPath + ';%INSTALL_DIR%'; [Environment]::SetEnvironmentVariable('Path', $newPath, 'User'); Write-Host '       Added to user PATH' -ForegroundColor Green; Write-Host ''; Write-Host '       Note: Please restart your terminal' -ForegroundColor Yellow } else { Write-Host '       Already in PATH' -ForegroundColor Green }"

echo.
echo ================================================
echo      Installation Successful!
echo ================================================
echo.
echo Install location: %INSTALL_DIR%\%BINARY_NAME%
echo Run command: cnfast --help
echo.
echo Note: If cnfast command not found, please restart your terminal
echo.

REM Verify installation
echo [5/5] Verifying installation...
"%INSTALL_DIR%\%BINARY_NAME%" --version 2>nul
if errorlevel 1 (
    echo       Please restart terminal to use cnfast command
) else (
    echo.
    echo       Installation verified successfully!
)

echo.
echo To uninstall, delete directory: %INSTALL_DIR%
echo.

pause