package main

import (
	"fmt"
	"runtime"
	"time"
)

// 打印版本信息
func printVersion() {
	fmt.Printf("%s version %s\n", name, version)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// 打印帮助信息
func printHelp() {
	fmt.Printf(`%s - 高性能命令行工具

用法:
    %s [选项] <命令> [参数]

可用命令:
    version        显示版本信息
    help           显示帮助信息
    info           显示系统信息
    benchmark      运行性能测试
    config         配置管理 (init|show)
    hash           计算文件哈希 <文件> <类型:md5|sha256>
    url            测试URL连接 <URL>
    find           查找文件 <目录> <模式>
    analyze        分析目录 <目录>
    process        并发处理文件 <目录> <工作线程数>
    grep           在文件中搜索文本 <模式> [文件/目录]
    interactive    交互模式

选项:
    -version       显示版本信息
    -help          显示帮助信息

示例:
    %s version
    %s info
    %s benchmark
    %s config init
    %s hash example.txt md5
    %s url https://github.com
    %s find . ".go"
    %s analyze /tmp
    %s process . 4
    %s grep "func main" .
    %s interactive

`, name, name, name, name, name, name, name, name, name, name, name, name, name)
}

// 打印系统信息
func printSystemInfo() {
	fmt.Println("系统信息:")
	fmt.Printf("  操作系统: %s\n", runtime.GOOS)
	fmt.Printf("  架构: %s\n", runtime.GOARCH)
	fmt.Printf("  Go版本: %s\n", runtime.Version())
	fmt.Printf("  CPU核心数: %d\n", runtime.NumCPU())
	fmt.Printf("  Goroutine数: %d\n", runtime.NumGoroutine())
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  内存使用: %.2f MB\n", float64(m.Alloc)/1024/1024)
}

// 运行性能测试
func runBenchmark() {
	fmt.Println("运行性能测试...")
	
	// 简单的性能测试示例
	n := 1000000
	fmt.Printf("测试: 执行 %d 次循环\n", n)
	
	start := time.Now()
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	end := time.Now()
	
	duration := end.Sub(start)
	fmt.Printf("结果: %d\n", sum)
	fmt.Printf("耗时: %v\n", duration)
	fmt.Printf("每次操作: %.2f ns\n", float64(duration.Nanoseconds())/float64(n))
}

// 处理基本命令
func handleBasicCommands(subcommand string) bool {
	switch subcommand {
	case "version":
		printVersion()
		return true
	case "help":
		printHelp()
		return true
	case "info":
		printSystemInfo()
		return true
	case "benchmark":
		runBenchmark()
		return true
	default:
		return false
	}
} 