@echo off
chcp 65001 >nul 2>&1
cls

echo ============================================
echo   OverstepLab - 靶场服务启动脚本
echo ============================================
echo.

:: 获取脚本所在目录
set "SCRIPT_DIR=%~dp0"
set "FRONTEND_DIR=%SCRIPT_DIR%src\frontend"
set "BACKEND_DIR=%SCRIPT_DIR%src\backend"
set "WEB_DIST_DIR=%BACKEND_DIR%\internal\web\dist"
set "EXE_NAME=oversteplab.exe"
set "EXE_PATH=%BACKEND_DIR%\%EXE_NAME%"

:: 检测 Go 环境
echo [1/5] 正在检测 Go 编译环境...
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到 Go 环境，请安装 Go 1.25+ 并添加到 PATH。
    echo 下载地址: https://go.dev/dl/
    pause
    exit /b 1
)
for /f "tokens=3" %%a in ('go version') do (
    echo [信息] 检测到 Go 版本: %%a
)

:: 检测 Node.js / pnpm 环境
echo.
echo [2/5] 正在检测前端构建环境...
set "HAS_PNPM=0"
set "HAS_NPM=0"
where pnpm >nul 2>nul
if not errorlevel 1 set "HAS_PNPM=1"
where npm >nul 2>nul
if not errorlevel 1 set "HAS_NPM=1"

if "%HAS_PNPM%"=="0" if "%HAS_NPM%"=="0" (
    echo [错误] 未检测到 pnpm 或 npm，请安装 Node.js 和 pnpm。
    echo 下载地址: https://nodejs.org/  安装后执行: npm install -g pnpm
    pause
    exit /b 1
)

:: 编译前端
echo.
echo [3/5] 正在编译前端...
cd /d "%FRONTEND_DIR%" 2>nul
if errorlevel 1 (
    echo [错误] 无法进入前端目录: %FRONTEND_DIR%
    pause
    exit /b 1
)

if not exist "node_modules" (
    echo [信息] 首次运行，正在安装前端依赖...
    if "%HAS_PNPM%"=="1" (
        pnpm install
    ) else (
        npm install
    )
    if errorlevel 1 (
        echo [错误] 前端依赖安装失败。
        pause
        exit /b 1
    )
)

if "%HAS_PNPM%"=="1" (
    call pnpm build
) else (
    call npm run build
)
if %errorlevel% neq 0 (
    echo [错误] 前端编译失败。
    pause
    exit /b 1
)

:: 复制前端产物到后端嵌入目录
echo.
echo [信息] 正在复制前端产物到后端嵌入目录...
if not exist "%WEB_DIST_DIR%" mkdir "%WEB_DIST_DIR%"
xcopy /e /y /i "%FRONTEND_DIR%\dist\*" "%WEB_DIST_DIR%\" >nul 2>&1
if errorlevel 4 (
    echo [警告] 复制前端产物到后端目录失败，尝试继续...
)
echo [信息] 前端编译完成。

:: 进入后端目录
cd /d "%BACKEND_DIR%" 2>nul
if errorlevel 1 (
    echo [错误] 无法进入后端目录: %BACKEND_DIR%
    pause
    exit /b 1
)
echo [信息] 当前工作目录: %CD%

:: 编译后端
echo.
echo [4/5] 正在编译后端服务...
if exist "%EXE_NAME%" (
    del /f /q "%EXE_NAME%" >nul 2>&1
)

go build -ldflags "-s -w" -o "%EXE_NAME%" ./cmd/server/
if %errorlevel% neq 0 (
    echo [错误] 编译失败，请检查代码或依赖。
    pause
    exit /b 1
)
echo [信息] 编译成功: %EXE_NAME%

:: 检查端口占用
netstat -ano | findstr :8080 >nul 2>&1
if not errorlevel 1 (
    echo [警告] 端口 8080 已被占用，尝试终止占用进程...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8080') do (
        taskkill /f /pid %%a >nul 2>&1
    )
    ping -n 2 127.0.0.1 >nul 2>&1
)

:: 启动服务
echo.
echo [5/5] 正在启动 OverstepLab 服务...
set "GIN_MODE=release"
echo [信息] 默认访问地址: http://localhost:8080
echo [信息] 按 Ctrl+C 停止服务
echo.
echo ============================================
echo   靶场已启动！请在浏览器中访问:
echo   http://localhost:8080
echo ============================================
echo.

"%EXE_NAME%"
if %errorlevel% neq 0 (
    echo [错误] 服务启动失败，可能是端口被占用或程序异常退出。
    echo [信息] 退出码: %errorlevel%
    pause
    exit /b 1
)

:: 服务停止后的提示
echo.
echo ============================================
echo [信息] 服务已停止。
echo ============================================
pause
