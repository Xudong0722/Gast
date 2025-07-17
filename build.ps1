# PowerShell Build Script for Gast
param(
    [string]$Target = "build",
    [switch]$Release,
    [switch]$Clean,
    [switch]$Help
)

Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "    Gast Build Script for Windows" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host ""

# 显示帮助信息
if ($Help) {
    Write-Host "Usage: .\build.ps1 [-Target <target>] [-Release] [-Clean] [-Help]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Targets:" -ForegroundColor Green
    Write-Host "  build      - Build for current platform (default)" -ForegroundColor White
    Write-Host "  release    - Build for all platforms" -ForegroundColor White
    Write-Host "  test       - Run tests" -ForegroundColor White
    Write-Host "  fmt        - Format Go code" -ForegroundColor White
    Write-Host "  vet        - Run go vet" -ForegroundColor White
    Write-Host "  clean      - Clean build artifacts" -ForegroundColor White
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Green
    Write-Host "  -Release   - Build optimized release version" -ForegroundColor White
    Write-Host "  -Clean     - Clean before building" -ForegroundColor White
    Write-Host "  -Help      - Show this help message" -ForegroundColor White
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Blue
    Write-Host "  .\build.ps1                    # Basic build" -ForegroundColor Gray
    Write-Host "  .\build.ps1 -Target release    # Build for all platforms" -ForegroundColor Gray
    Write-Host "  .\build.ps1 -Clean             # Clean and build" -ForegroundColor Gray
    Write-Host "  .\build.ps1 -Target test       # Run tests" -ForegroundColor Gray
    exit 0
}

# 检查Go是否安装
Write-Host "[INFO] Checking Go installation..." -ForegroundColor Blue
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Go not found"
    }
    Write-Host "[SUCCESS] Found Go: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Go compiler not found!" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/ and add it to PATH." -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# 清理函数
function Clean-Build {
    Write-Host "[INFO] Cleaning build artifacts..." -ForegroundColor Yellow
    if (Test-Path "build") {
        Remove-Item -Recurse -Force "build"
        Write-Host "[INFO] Removed build directory" -ForegroundColor Blue
    }
    go clean 2>$null
    Write-Host "[SUCCESS] Clean completed" -ForegroundColor Green
}

# 构建函数
function Build-Current {
    Write-Host "[INFO] Building Gast for Windows..." -ForegroundColor Yellow
    
    # 创建构建目录
    if (!(Test-Path "build")) {
        New-Item -ItemType Directory -Path "build" | Out-Null
        Write-Host "[INFO] Created build directory" -ForegroundColor Blue
    }
    
    # 设置构建参数
    $buildArgs = @("build", "-o", "build\gast.exe")
    if ($Release) {
        $buildArgs += @("-ldflags", "-s -w")  # 压缩二进制文件
        Write-Host "[INFO] Building optimized release version..." -ForegroundColor Cyan
    }
    $buildArgs += "."
    
    # 执行构建
    & go @buildArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "[SUCCESS] Build completed successfully!" -ForegroundColor Green
        Write-Host "Binary created: build\gast.exe" -ForegroundColor Blue
        
        # 获取文件大小
        $fileSize = (Get-Item "build\gast.exe").Length
        $fileSizeMB = [math]::Round($fileSize / 1MB, 2)
        Write-Host "Binary size: $fileSizeMB MB" -ForegroundColor Blue
        
        Write-Host ""
        Write-Host "[INFO] Testing the binary:" -ForegroundColor Blue
        & ".\build\gast.exe" version
        Write-Host ""
        Write-Host "[INFO] You can now use: .\build\gast.exe [command]" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] Build failed!" -ForegroundColor Red
        exit 1
    }
}

# 发布构建函数
function Build-Release {
    Write-Host "[INFO] Building release versions for all platforms..." -ForegroundColor Yellow
    
    if (!(Test-Path "build\release")) {
        New-Item -ItemType Directory -Path "build\release" -Force | Out-Null
    }
    
    $platforms = @(
        @{GOOS="windows"; GOARCH="amd64"; NAME="gast-windows-amd64.exe"},
        @{GOOS="linux"; GOARCH="amd64"; NAME="gast-linux-amd64"},
        @{GOOS="linux"; GOARCH="arm64"; NAME="gast-linux-arm64"},
        @{GOOS="darwin"; GOARCH="amd64"; NAME="gast-darwin-amd64"},
        @{GOOS="darwin"; GOARCH="arm64"; NAME="gast-darwin-arm64"}
    )
    
    foreach ($platform in $platforms) {
        Write-Host "Building for $($platform.GOOS)/$($platform.GOARCH)..." -ForegroundColor Cyan
        
        $env:GOOS = $platform.GOOS
        $env:GOARCH = $platform.GOARCH
        
        $buildArgs = @("build", "-ldflags", "-s -w", "-o", "build\release\$($platform.NAME)", ".")
        & go @buildArgs
        
        if ($LASTEXITCODE -eq 0) {
            $size = (Get-Item "build\release\$($platform.NAME)").Length
            $sizeMB = [math]::Round($size / 1MB, 2)
            Write-Host "  ✓ $($platform.NAME) ($sizeMB MB)" -ForegroundColor Green
        } else {
            Write-Host "  ✗ Failed to build for $($platform.GOOS)/$($platform.GOARCH)" -ForegroundColor Red
        }
    }
    
    # 重置环境变量
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    
    Write-Host ""
    Write-Host "[SUCCESS] Release build completed!" -ForegroundColor Green
    Write-Host "Binaries created in: build\release\" -ForegroundColor Blue
    
    # 显示生成的文件
    Get-ChildItem "build\release\" | ForEach-Object {
        $sizeMB = [math]::Round($_.Length / 1MB, 2)
        Write-Host "  - $($_.Name) ($sizeMB MB)" -ForegroundColor Gray
    }
}

# 执行清理（如果需要）
if ($Clean) {
    Clean-Build
    Write-Host ""
}

# 根据目标执行相应操作
switch ($Target.ToLower()) {
    "build" {
        Build-Current
    }
    
    "release" {
        Build-Release
    }
    
    "test" {
        Write-Host "[INFO] Running tests..." -ForegroundColor Yellow
        go test -v ./...
        if ($LASTEXITCODE -eq 0) {
            Write-Host "[SUCCESS] All tests passed!" -ForegroundColor Green
        } else {
            Write-Host "[ERROR] Tests failed!" -ForegroundColor Red
            exit 1
        }
    }
    
    "fmt" {
        Write-Host "[INFO] Formatting Go code..." -ForegroundColor Yellow
        go fmt ./...
        Write-Host "[SUCCESS] Code formatting completed!" -ForegroundColor Green
    }
    
    "vet" {
        Write-Host "[INFO] Running go vet..." -ForegroundColor Yellow
        go vet ./...
        if ($LASTEXITCODE -eq 0) {
            Write-Host "[SUCCESS] Go vet completed successfully!" -ForegroundColor Green
        } else {
            Write-Host "[ERROR] Go vet found issues!" -ForegroundColor Red
            exit 1
        }
    }
    
    "clean" {
        Clean-Build
    }
    
    default {
        Write-Host "[ERROR] Unknown target: $Target" -ForegroundColor Red
        Write-Host "Use -Help to see available targets" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host ""
Write-Host "Done!" -ForegroundColor Green 