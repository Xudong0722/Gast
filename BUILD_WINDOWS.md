# Windows 构建指南

本文档说明如何在Windows系统上构建Gast项目。

## 前提条件

1. **安装Go语言**
   - 从 [https://golang.org/dl/](https://golang.org/dl/) 下载Go安装包
   - 安装后确保`go`命令在PATH中可用
   - 验证安装：打开命令提示符或PowerShell，运行 `go version`

2. **安装Git** (可选，用于克隆仓库)
   - 从 [https://git-scm.com/download/win](https://git-scm.com/download/win) 下载安装

3. **安装CMake** (可选，如果使用CMake构建)
   - 从 [https://cmake.org/download/](https://cmake.org/download/) 下载安装
   - 安装时选择"Add CMake to system PATH"

## 构建方法

### 方法1: 使用PowerShell脚本 (推荐)

这是最简单的方法，支持多种构建选项：

```powershell
# 基本构建
.\build.ps1

# 显示帮助
.\build.ps1 -Help

# 优化发布版本
.\build.ps1 -Release

# 清理并构建
.\build.ps1 -Clean

# 构建所有平台版本
.\build.ps1 -Target release

# 运行测试
.\build.ps1 -Target test

# 格式化代码
.\build.ps1 -Target fmt
```

### 方法2: 使用批处理脚本

双击`build.bat`文件或在命令提示符中运行：

```cmd
build.bat
```

### 方法3: 使用CMake

```cmd
# 创建构建目录
mkdir cmake-build
cd cmake-build

# 配置项目
cmake ..

# 编译
cmake --build .

# 清理
cmake --build . --target clean-go

# 构建所有平台版本
cmake --build . --target build-all
```

### 方法4: 直接使用Go命令

```cmd
# 基本构建
go build -o build\gast.exe .

# 优化构建
go build -ldflags "-s -w" -o build\gast.exe .

# 清理
go clean
```

## 交叉编译

如果需要为其他平台构建：

```cmd
# Linux AMD64
set GOOS=linux
set GOARCH=amd64
go build -o build\gast-linux-amd64 .

# macOS AMD64
set GOOS=darwin
set GOARCH=amd64
go build -o build\gast-darwin-amd64 .

# 重置环境变量
set GOOS=
set GOARCH=
```

或使用PowerShell：

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o build\gast-linux-amd64 .
```

## 常见问题

### 1. Go命令未找到

**错误**: `'go' is not recognized as an internal or external command`

**解决方案**: 
- 确认Go已正确安装
- 检查Go安装目录是否已添加到系统PATH中
- 重启命令提示符或PowerShell

### 2. PowerShell执行策略限制

**错误**: `cannot be loaded because running scripts is disabled on this system`

**解决方案**:
```powershell
# 以管理员身份运行PowerShell，然后执行：
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 3. 构建失败

**错误**: 编译时出现各种错误

**解决方案**:
- 确保所有Go源文件在同一目录
- 检查go.mod文件是否存在
- 运行 `go mod tidy` 检查依赖
- 查看错误信息并检查代码语法

### 4. CMake未找到

**错误**: `cmake is not recognized as an internal or external command`

**解决方案**:
- 安装CMake并添加到PATH
- 或者使用其他构建方法（PowerShell脚本或批处理）

## 输出文件

构建成功后，可执行文件将生成在以下位置：

- **PowerShell/批处理脚本**: `build\gast.exe`
- **CMake**: `cmake-build\bin\gast.exe`
- **发布版本**: `build\release\` 目录下的各平台版本

## 测试构建

构建完成后，可以测试二进制文件：

```cmd
# 显示版本信息
build\gast.exe version

# 显示帮助
build\gast.exe help

# 运行系统信息
build\gast.exe info
```

## 开发建议

1. **使用PowerShell脚本**进行日常开发，它提供了最完整的功能
2. **使用批处理脚本**进行快速构建
3. **使用CMake**进行复杂的构建配置
4. 在提交代码前运行 `.\build.ps1 -Target fmt` 格式化代码
5. 运行 `.\build.ps1 -Target test` 确保测试通过

## 集成开发环境

### Visual Studio Code
- 安装Go扩展
- 配置tasks.json使用构建脚本
- 使用集成终端运行构建命令

### GoLand/IntelliJ
- 导入Go模块
- 配置外部工具使用构建脚本
- 使用内置终端运行构建命令 