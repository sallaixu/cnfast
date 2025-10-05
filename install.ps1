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
Write-Host "     cnfast Windows 一键安装脚本" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# [1/5] 检测系统架构
Write-Host "[1/5] 检测系统架构..." -ForegroundColor Yellow
$Arch = $env:PROCESSOR_ARCHITECTURE
$ArchSuffix = ""

switch ($Arch) {
    "AMD64" {
        Write-Host "      检测到: AMD64" -ForegroundColor Green
        $ArchSuffix = "amd64"
    }
    "ARM64" {
        Write-Host "      检测到: ARM64" -ForegroundColor Green
        $ArchSuffix = "arm64"
    }
    "x86" {
        Write-Host "      检测到: x86" -ForegroundColor Green
        $ArchSuffix = "386"
    }
    default {
        Write-Host "      错误: 不支持的架构 $Arch" -ForegroundColor Red
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
    # [2/5] 下载 cnfast
    Write-Host "[2/5] 下载 cnfast..." -ForegroundColor Yellow
    Write-Host ""
    
    # 设置 TLS 1.2
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    
    # 优先使用 BitsTransfer（带进度）
    if (Get-Command Start-BitsTransfer -ErrorAction SilentlyContinue) {
        Write-Host "      使用 BITS 传输服务（带进度显示）..." -ForegroundColor Cyan
        try {
            Import-Module BitsTransfer
            Start-BitsTransfer -Source $DownloadUrl -Destination $TmpFile -DisplayName "下载 cnfast" -Description "正在下载，请稍候..."
            Write-Host "      下载完成！" -ForegroundColor Green
        }
        catch {
            Write-Host "      BITS 下载失败，尝试备用方式..." -ForegroundColor Yellow
            throw
        }
    }
    else {
        # 备用方案：使用 WebClient
        Write-Host "      下载中..." -ForegroundColor Cyan
        $webClient = New-Object System.Net.WebClient
        
        # 注册进度事件
        $progressRegistered = $false
        try {
            Register-ObjectEvent -InputObject $webClient -EventName DownloadProgressChanged -SourceIdentifier WebClient.DownloadProgressChanged -Action {
                $percent = $EventArgs.ProgressPercentage
                $received = $EventArgs.BytesReceived / 1MB
                $total = $EventArgs.TotalBytesToReceive / 1MB
                Write-Progress -Activity "下载 cnfast" -Status "$([Math]::Round($received, 2)) MB / $([Math]::Round($total, 2)) MB" -PercentComplete $percent
            } | Out-Null
            $progressRegistered = $true
        }
        catch {
            # 如果注册事件失败，继续不显示进度
        }
        
        try {
            $webClient.DownloadFile($DownloadUrl, $TmpFile)
            if ($progressRegistered) {
                Write-Progress -Activity "下载 cnfast" -Completed
            }
            Write-Host "      下载完成！" -ForegroundColor Green
        }
        finally {
            if ($progressRegistered) {
                Unregister-Event -SourceIdentifier WebClient.DownloadProgressChanged -ErrorAction SilentlyContinue
            }
            $webClient.Dispose()
        }
    }
    
    Write-Host ""
    
    # 检查文件是否存在且大小正常
    if (-not (Test-Path $TmpFile)) {
        Write-Host "      错误: 下载失败，文件不存在" -ForegroundColor Red
        exit 1
    }
    
    $fileSize = (Get-Item $TmpFile).Length
    if ($fileSize -lt 1024) {
        Write-Host "      错误: 下载文件异常（小于 1KB）" -ForegroundColor Red
        exit 1
    }
    
    # [3/5] 安装到本地
    Write-Host "[3/5] 安装到本地..." -ForegroundColor Yellow
    
    # 创建安装目录
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # 如果目标文件存在，先删除
    $TargetFile = Join-Path $InstallDir $BinaryName
    if (Test-Path $TargetFile) {
        try {
            Remove-Item $TargetFile -Force
        }
        catch {
            Write-Host "      警告: 无法删除旧版本文件，可能正在使用中" -ForegroundColor Yellow
            Write-Host "      请关闭所有 cnfast 进程后重试" -ForegroundColor Yellow
            exit 1
        }
    }
    
    # 移动文件到安装目录
    Move-Item -Path $TmpFile -Destination $TargetFile -Force
    Write-Host "      安装完成！" -ForegroundColor Green
    Write-Host "      位置: $TargetFile" -ForegroundColor Cyan
    Write-Host ""
    
    # [4/5] 配置环境变量
    Write-Host "[4/5] 配置环境变量..." -ForegroundColor Yellow
    
    # 获取当前用户的 PATH
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    # 检查是否已经在 PATH 中
    if ($UserPath -notlike "*$InstallDir*") {
        $NewPath = $UserPath + ";" + $InstallDir
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-Host "      已添加到用户 PATH" -ForegroundColor Green
        
        # 更新当前会话的 PATH
        $env:Path = $env:Path + ";" + $InstallDir
        
        Write-Host ""
        Write-Host "      提示: 如果当前终端无法使用 cnfast 命令，" -ForegroundColor Yellow
        Write-Host "            请重新打开一个新的终端窗口" -ForegroundColor Yellow
    }
    else {
        Write-Host "      安装目录已在 PATH 中" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host "     安装成功！" -ForegroundColor Green
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "安装位置: $TargetFile" -ForegroundColor Cyan
    Write-Host "运行命令: cnfast --help" -ForegroundColor Cyan
    Write-Host ""
    
    # [5/5] 验证安装
    Write-Host "[5/5] 验证安装..." -ForegroundColor Yellow
    try {
        $version = & $TargetFile --version 2>&1
        Write-Host $version
        Write-Host ""
        Write-Host "      ✓ 安装验证成功！" -ForegroundColor Green
    }
    catch {
        Write-Host "      安装完成，请重新打开终端后使用 cnfast 命令" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "卸载方法: Remove-Item '$InstallDir' -Recurse -Force" -ForegroundColor Cyan
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "安装失败: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}
finally {
    # 清理临时文件
    if (Test-Path $TmpFile) {
        Remove-Item -Path $TmpFile -Force -ErrorAction SilentlyContinue
    }
}