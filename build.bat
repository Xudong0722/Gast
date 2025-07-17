@echo off
setlocal EnableDelayedExpansion

echo =======================================
echo     Gast Build Script for Windows
echo =======================================
echo.

REM 检查Go是否安装
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Go compiler not found. 
    echo Please install Go from https://golang.org/dl/ and add it to PATH.
    echo.
    pause
    exit /b 1
)

REM 显示Go版本
echo [INFO] Found Go:
go version
echo.

REM 创建构建目录
if not exist "build" (
    echo [INFO] Creating build directory...
    mkdir build
)

REM 编译程序
echo [INFO] Compiling Gast...
go build -o build\gast.exe .

if %errorlevel% equ 0 (
    echo.
    echo [SUCCESS] Build completed successfully!
    echo Binary created: build\gast.exe
    echo.
    echo [INFO] Testing the binary:
    build\gast.exe version
    echo.
    echo [INFO] You can now use: build\gast.exe [command]
) else (
    echo.
    echo [ERROR] Build failed!
    echo Please check the error messages above.
    echo.
    pause
    exit /b 1
)

echo.
echo Press any key to exit...
pause >nul 