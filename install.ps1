# cnfast Windows 一键安装脚本
# PowerShell 脚本，支持 Windows 10/11

# 设置控制台编码为 UTF-8（关键：解决中文乱码）
$PSDefaultParameterValues['Out-File:Encoding'] = 'utf8'
$OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# 设置错误处理
$ErrorActionPreference = "Stop"

# 定义变量
$BaseUrl = "https://gitee.com/sallai/cnfast/releases/download/latest"
$BinaryName = "cnfast.exe"
$InstallDir = "$env:LOCALAPPDATA\cnfast"

Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "     cnfast Windows Installer" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# [1/5] Detect architecture
Write-Host "[1/5] Detecting system architecture..." -ForegroundColor Yellow
$Arch = $env:PROCESSOR_ARCHITECTURE
$ArchSuffix = ""

switch ($Arch) {
    "AMD64" {
        Write-Host "      Detected: AMD64" -ForegroundColor Green
        $ArchSuffix = "amd64"
    }
    "ARM64" {
        Write-Host "      Detected: ARM64" -ForegroundColor Green
        $ArchSuffix = "arm64"
    }
    "x86" {
        Write-Host "      Detected: x86" -ForegroundColor Green
        $ArchSuffix = "386"
    }
    default {
        Write-Host "      Error: Unsupported architecture $Arch" -ForegroundColor Red
        exit 1
    }
}

# 构建下载URL
$DownloadUrl = "$BaseUrl/cnfast-windows-$ArchSuffix.exe"
Write-Host "      URL: $DownloadUrl" -ForegroundColor Cyan
Write-Host ""

# 创建临时文件路径
$TmpFile = Join-Path $env:TEMP "cnfast_$(Get-Random).exe"

try {
    # [2/5] Download cnfast
    Write-Host "[2/5] Downloading cnfast..." -ForegroundColor Yellow
    Write-Host ""
    
    # Set TLS 1.2
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    
    # Try BitsTransfer first (with progress)
    if (Get-Command Start-BitsTransfer -ErrorAction SilentlyContinue) {
        Write-Host "      Using BITS transfer (with progress)..." -ForegroundColor Cyan
        try {
            Import-Module BitsTransfer
            Start-BitsTransfer -Source $DownloadUrl -Destination $TmpFile -DisplayName "Downloading cnfast" -Description "Please wait..."
            Write-Host "      Download completed!" -ForegroundColor Green
        }
        catch {
            Write-Host "      BITS failed, trying alternative..." -ForegroundColor Yellow
            throw
        }
    }
    else {
        # Fallback: WebClient
        Write-Host "      Downloading..." -ForegroundColor Cyan
        $webClient = New-Object System.Net.WebClient
        
        # Register progress event
        $progressRegistered = $false
        try {
            Register-ObjectEvent -InputObject $webClient -EventName DownloadProgressChanged -SourceIdentifier WebClient.DownloadProgressChanged -Action {
                $percent = $EventArgs.ProgressPercentage
                $received = $EventArgs.BytesReceived / 1MB
                $total = $EventArgs.TotalBytesToReceive / 1MB
                Write-Progress -Activity "Downloading cnfast" -Status "$([Math]::Round($received, 2)) MB / $([Math]::Round($total, 2)) MB" -PercentComplete $percent
            } | Out-Null
            $progressRegistered = $true
        }
        catch {
            # If event registration fails, continue without progress
        }
        
        try {
            $webClient.DownloadFile($DownloadUrl, $TmpFile)
            if ($progressRegistered) {
                Write-Progress -Activity "Downloading cnfast" -Completed
            }
            Write-Host "      Download completed!" -ForegroundColor Green
        }
        finally {
            if ($progressRegistered) {
                Unregister-Event -SourceIdentifier WebClient.DownloadProgressChanged -ErrorAction SilentlyContinue
            }
            $webClient.Dispose()
        }
    }
    
    Write-Host ""
    
    # Check file
    if (-not (Test-Path $TmpFile)) {
        Write-Host "      Error: Download failed" -ForegroundColor Red
        exit 1
    }
    
    $fileSize = (Get-Item $TmpFile).Length
    if ($fileSize -lt 1024) {
        Write-Host "      Error: File too small (< 1KB)" -ForegroundColor Red
        exit 1
    }
    
    # [3/5] Install
    Write-Host "[3/5] Installing..." -ForegroundColor Yellow
    
    # Create install directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # Remove old version if exists
    $TargetFile = Join-Path $InstallDir $BinaryName
    if (Test-Path $TargetFile) {
        try {
            Remove-Item $TargetFile -Force
        }
        catch {
            Write-Host "      Warning: Cannot delete old version (may be in use)" -ForegroundColor Yellow
            Write-Host "      Please close all cnfast processes and try again" -ForegroundColor Yellow
            exit 1
        }
    }
    
    # Move file
    Move-Item -Path $TmpFile -Destination $TargetFile -Force
    Write-Host "      Installation completed!" -ForegroundColor Green
    Write-Host "      Location: $TargetFile" -ForegroundColor Cyan
    Write-Host ""
    
    # [4/5] Configure PATH
    Write-Host "[4/5] Configuring environment variables..." -ForegroundColor Yellow
    
    # Get user PATH
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    # Check if already in PATH
    if ($UserPath -notlike "*$InstallDir*") {
        $NewPath = $UserPath + ";" + $InstallDir
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Host "      Added to user PATH" -ForegroundColor Green
        
        # Update current session
        $env:Path = $env:Path + ";" + $InstallDir
        
        Write-Host ""
        Write-Host "      Note: If command not found," -ForegroundColor Yellow
        Write-Host "            please restart your terminal" -ForegroundColor Yellow
    }
    else {
        Write-Host "      Already in PATH" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host "     Installation Successful!" -ForegroundColor Green
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Install location: $TargetFile" -ForegroundColor Cyan
    Write-Host "Run command: cnfast --help" -ForegroundColor Cyan
    Write-Host ""
    
    # [5/5] Verify
    Write-Host "[5/5] Verifying installation..." -ForegroundColor Yellow
    try {
        $version = & $TargetFile --version 2>&1
        Write-Host $version
        Write-Host ""
        Write-Host "      Installation verified!" -ForegroundColor Green
    }
    catch {
        Write-Host "      Please restart terminal to use cnfast" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "Uninstall: Remove-Item '$InstallDir' -Recurse -Force" -ForegroundColor Cyan
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "Installation failed: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}
finally {
    # 清理临时文件
    if (Test-Path $TmpFile) {
        Remove-Item -Path $TmpFile -Force -ErrorAction SilentlyContinue
    }
}