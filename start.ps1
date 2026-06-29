#!/usr/bin/env pwsh
#Requires -Version 5.1

[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

$Host.UI.RawUI.WindowTitle = "OverstepLab - 靶场服务启动脚本"

Clear-Host

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "   OverstepLab - 靶场服务启动脚本" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# 获取脚本所在目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$FrontendDir = Join-Path $ScriptDir "src\frontend"
$BackendDir = Join-Path $ScriptDir "src\backend"
$WebDistDir = Join-Path $BackendDir "internal\web\dist"
$ExeName = "oversteplab.exe"
$ExePath = Join-Path $BackendDir $ExeName

# 检测 Go 环境
Write-Host "[1/5] 正在检测 Go 编译环境..." -ForegroundColor Yellow
$GoVersion = $null
try {
    $GoVersion = (go version 2>$null) | Select-Object -First 1
} catch {
    $GoVersion = $null
}

if (-not $GoVersion) {
    Write-Host "[错误] 未检测到 Go 环境，请安装 Go 1.25+ 并添加到 PATH。" -ForegroundColor Red
    Write-Host "下载地址: https://go.dev/dl/" -ForegroundColor Red
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

Write-Host "[信息] 检测到 Go 版本: $GoVersion" -ForegroundColor Green

# 检测 Node.js / pnpm 环境
Write-Host ""
Write-Host "[2/5] 正在检测前端构建环境..." -ForegroundColor Yellow
$HasPnpm = $null -ne (Get-Command pnpm -ErrorAction SilentlyContinue)
$HasNpm = $null -ne (Get-Command npm -ErrorAction SilentlyContinue)

if (-not $HasPnpm -and -not $HasNpm) {
    Write-Host "[错误] 未检测到 pnpm 或 npm，请安装 Node.js 和 pnpm。" -ForegroundColor Red
    Write-Host "下载地址: https://nodejs.org/  安装后执行: npm install -g pnpm" -ForegroundColor Red
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

# 编译前端
Write-Host ""
Write-Host "[3/5] 正在编译前端..." -ForegroundColor Yellow
if (-not (Test-Path $FrontendDir)) {
    Write-Host "[错误] 无法找到前端目录: $FrontendDir" -ForegroundColor Red
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

Push-Location $FrontendDir

if (-not (Test-Path "node_modules")) {
    Write-Host "[信息] 首次运行，正在安装前端依赖..." -ForegroundColor Cyan
    try {
        if ($HasPnpm) {
            pnpm install 2>&1 | ForEach-Object { Write-Host $_ -ForegroundColor Gray }
        } else {
            npm install 2>&1 | ForEach-Object { Write-Host $_ -ForegroundColor Gray }
        }
        if ($LASTEXITCODE -ne 0) { throw "依赖安装失败" }
    } catch {
        Write-Host "[错误] 前端依赖安装失败。" -ForegroundColor Red
        Pop-Location
        Write-Host ""
        Write-Host "按任意键退出..." -ForegroundColor Gray
        $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
        exit 1
    }
}

try {
    if ($HasPnpm) {
        pnpm build 2>&1 | ForEach-Object { Write-Host $_ -ForegroundColor Gray }
    } else {
        npm run build 2>&1 | ForEach-Object { Write-Host $_ -ForegroundColor Gray }
    }
    if ($LASTEXITCODE -ne 0) { throw "前端编译失败" }
} catch {
    Write-Host "[错误] 前端编译失败。" -ForegroundColor Red
    Pop-Location
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

# 复制前端产物到后端嵌入目录
if (-not (Test-Path $WebDistDir)) {
    New-Item -ItemType Directory -Path $WebDistDir -Force | Out-Null
}
$FrontendDist = Join-Path $FrontendDir "dist"
if (Test-Path $FrontendDist) {
    try {
        Copy-Item -Path "$FrontendDist\*" -Destination $WebDistDir -Recurse -Force -ErrorAction Stop
    } catch {
        Write-Host "[警告] 复制前端产物到后端目录失败，尝试继续..." -ForegroundColor Yellow
    }
}

Write-Host "[信息] 前端编译完成。" -ForegroundColor Green
Pop-Location

# 进入后端目录
if (-not (Test-Path $BackendDir)) {
    Write-Host "[错误] 无法找到后端目录: $BackendDir" -ForegroundColor Red
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

Push-Location $BackendDir

# 编译后端
Write-Host ""
Write-Host "[4/5] 正在编译后端服务..." -ForegroundColor Yellow
if (Test-Path $ExeName) {
    Remove-Item $ExeName -Force -ErrorAction SilentlyContinue
}

try {
    go build -ldflags "-s -w" -o $ExeName ./cmd/server/ 2>&1 | ForEach-Object {
        Write-Host $_ -ForegroundColor Gray
    }
    if ($LASTEXITCODE -ne 0) {
        throw "编译失败"
    }
} catch {
    Write-Host "[错误] 编译失败，请检查代码或依赖。" -ForegroundColor Red
    Write-Host "详细错误: $_" -ForegroundColor Red
    Pop-Location
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

Write-Host "[信息] 编译成功: $ExeName" -ForegroundColor Green

# 启动服务
Write-Host ""
Write-Host "[5/5] 正在启动 OverstepLab 服务..." -ForegroundColor Yellow
Write-Host "[信息] 默认访问地址: http://localhost:5000" -ForegroundColor Green
Write-Host "[信息] 按 Ctrl+C 停止服务" -ForegroundColor Green
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  靶场已启动！请在浏览器中访问:" -ForegroundColor Cyan
Write-Host "  http://localhost:5000" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

$env:OVERSTEPLAB_PORT = "5000"
$env:GIN_MODE = "release"

# 启动进程并等待
try {
    $Process = Start-Process -FilePath $ExePath -NoNewWindow -Wait -PassThru
    $ExitCode = $Process.ExitCode
} catch {
    Write-Host "[错误] 启动服务失败: $_" -ForegroundColor Red
    Pop-Location
    Write-Host ""
    Write-Host "按任意键退出..." -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

Pop-Location

# 服务停止后的提示
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "[信息] 服务已停止 (退出码: $ExitCode)。" -ForegroundColor Yellow
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
