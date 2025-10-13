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
echo.

REM Use PowerShell BitsTransfer or curl for progress
where /q bitsadmin
if %errorlevel% equ 0 (
    echo       Using BitsTransfer for download with progress...
    powershell -NoProfile -ExecutionPolicy Bypass -Command "Import-Module BitsTransfer; Start-BitsTransfer -Source '%DOWNLOAD_URL%' -Destination '%TMP_FILE%' -DisplayName 'Downloading cnfast' -Description 'Please wait...'; Write-Host '       Download completed!' -ForegroundColor Green"
) else (
    echo       Downloading... Please wait...
    powershell -NoProfile -ExecutionPolicy Bypass -Command "[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; $wc = New-Object System.Net.WebClient; $wc.DownloadFile('%DOWNLOAD_URL%', '%TMP_FILE%'); Write-Host '       Download completed!' -ForegroundColor Green"
)

if errorlevel 1 (
    echo       错误：下载失败，请检查你的网络连接
    exit /b 1
)

REM Check file exists
if not exist "%TMP_FILE%" (
    echo       错误：下载失败，文件不存在
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
        echo       警告：无法删除旧版本，可能正在使用
        echo       请关闭所有cnfast进程后重试
        del /f /q "%TMP_FILE%" 2>nul
        exit /b 1
    )
)

REM Move file to install directory
move /y "%TMP_FILE%" "%INSTALL_DIR%\%BINARY_NAME%" >nul

if errorlevel 1 (
    echo       错误：安装失败
    exit /b 1
)

echo       安装完成！
echo.

REM Add to PATH
echo [4/5] 配置环境变量中....
echo.

REM Use PowerShell to add to user PATH
powershell -NoProfile -ExecutionPolicy Bypass -Command "$userPath = [Environment]::GetEnvironmentVariable('Path', 'User'); if ($userPath -notlike '*%INSTALL_DIR%*') { $newPath = $userPath + ';%INSTALL_DIR%'; [Environment]::SetEnvironmentVariable('Path', $newPath, 'User'); Write-Host '       Added to user PATH' -ForegroundColor Green; Write-Host ''; Write-Host '       Note: Please restart your terminal' -ForegroundColor Yellow } else { Write-Host '       Already in PATH' -ForegroundColor Green }"

echo.
echo ================================================
echo      安装
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
    echo       提示：如果无法使用cnfast命令，请重新打开终端
) else (
    echo.
    echo       安装验证成功！
)

echo.
echo 卸载请删除此目录文件即可: %INSTALL_DIR%
echo.

pause