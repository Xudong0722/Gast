# Gast

一个用Go语言构建的高性能命令行工具，提供多种实用功能。

## 功能特性

- 🚀 **高性能** - 基于Go语言，启动快速，内存使用低
- 🔧 **多功能** - 集成文件操作、网络工具、系统信息等常用功能
- 💫 **并发处理** - 支持多线程并发文件处理
- 🎯 **配置管理** - 支持JSON配置文件
- 🔍 **文件查找** - 快速文件搜索和分析
- 📊 **性能测试** - 内置基准测试功能
- 🌐 **网络工具** - URL连接测试
- 🔐 **文件哈希** - MD5/SHA256哈希计算
- 🔍 **文本搜索** - 类似grep的强大文本搜索功能
- 📱 **交互模式** - 支持交互式命令行

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/Xudong0722/Gast.git
cd gast

# 编译
make build

# 安装到系统 (可选)
make install
```

### 直接编译

```bash
go build -o gast .
```

## 使用方法

### 基本命令

```bash
# 显示帮助
./gast help

# 显示版本信息
./gast version

# 显示系统信息
./gast info

# 运行性能测试
./gast benchmark
```

### 配置管理

```bash
# 初始化配置文件
./gast config init

# 查看当前配置
./gast config show
```

### 文件操作

```bash
# 计算文件哈希
./gast hash filename.txt md5
./gast hash filename.txt sha256

# 查找文件
./gast find /path/to/directory "*.go"

# 分析目录
./gast analyze /path/to/directory

# 并发处理文件 (使用4个工作线程)
./gast process /path/to/directory 4
```

### 网络工具

```bash
# 测试URL连接
./gast url https://github.com
./gast url google.com
```

### 文本搜索 (Grep)

```bash
# 基本搜索
./gast grep "pattern" file.txt

# 忽略大小写搜索
./gast grep -i "Hello" file.txt

# 显示行号
./gast grep -n "func main" main.go

# 递归搜索目录
./gast grep -r "TODO" src/

# 反向匹配（显示不匹配的行）
./gast grep -v "test" file.txt

# 只显示匹配行数
./gast grep -c "import" *.go

# 只显示匹配的文件名
./gast grep -l "fmt.Printf" *.go

# 组合选项
./gast grep -r -i -n "error" .
```

### 交互模式

```bash
# 进入交互模式
./gast interactive

# 在交互模式中可以使用以下命令:
gast> info
gast> version
gast> benchmark
gast> config
gast> quit
```

## 开发

### 构建

```bash
# 编译
make build

# 清理
make clean

# 格式化代码
make fmt

# 静态检查
make vet

# 运行测试
make test
```

### 发布

```bash
# 构建所有平台的发布版本
make release
```

这将在`build/release/`目录中生成以下文件：
- `gast-linux-amd64`
- `gast-linux-arm64`
- `gast-windows-amd64.exe`
- `gast-darwin-amd64`
- `gast-darwin-arm64`

### 版本控制

项目使用`.gitignore`文件忽略以下内容：
- `build/` 目录（编译产物）
- 二进制文件（`gast`, `*.exe`等）
- 测试和覆盖率文件
- IDE配置文件
- 操作系统生成的文件

## 配置文件

配置文件位于`~/.gast.json`，包含以下选项：

```json
{
  "log_level": "info",
  "output_dir": "./output",
  "max_workers": 4,
  "timeout": 30,
  "enable_color": true
}
```

## 项目结构

```
├── main.go                 # 主程序入口和命令路由
├── cmd_basic.go           # 基本命令 (version, help, info, benchmark)
├── cmd_config.go          # 配置管理命令
├── cmd_file.go            # 文件操作命令 (hash, find, analyze, process)
├── cmd_network.go         # 网络相关命令 (url)
├── cmd_grep.go            # 文本搜索命令
├── cmd_interactive.go     # 交互模式
├── config.go              # 配置文件处理
├── utils.go               # 工具函数和核心功能
├── Makefile               # 构建脚本
├── .gitignore             # Git忽略文件
├── ARCHITECTURE.md        # 架构说明文档
└── README.md              # 项目文档
```

## 性能

Gast采用Go语言开发，具有以下性能特征：

- 启动时间 < 10ms
- 内存使用 < 1MB
- 支持高并发文件处理
- 跨平台支持
- 模块化设计，易于扩展

## 扩展开发

### 添加新命令

项目采用模块化设计，添加新命令非常简单：

1. **创建命令文件**：在项目根目录创建 `cmd_xxx.go` 文件
2. **实现命令处理器**：
   ```go
   package main
   
   import "fmt"
   
   func handleMyCommand(args []string) {
       // 命令逻辑
   }
   
   func handleMyCommands(subcommand string, args []string) bool {
       switch subcommand {
       case "mycmd":
           handleMyCommand(args)
           return true
       default:
           return false
       }
   }
   ```

3. **注册命令**：在 `main.go` 的 `routeCommand` 函数中添加：
   ```go
   if handleMyCommands(subcommand, args) {
       return
   }
   ```

4. **更新帮助**：在 `cmd_basic.go` 的 `printHelp` 函数中添加命令说明

### 命令分类

- `cmd_basic.go` - 基础系统命令
- `cmd_config.go` - 配置管理
- `cmd_file.go` - 文件操作
- `cmd_network.go` - 网络工具
- `cmd_grep.go` - 文本搜索
- `cmd_interactive.go` - 交互模式

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

本项目采用MIT许可证 - 详见[LICENSE](LICENSE)文件。
