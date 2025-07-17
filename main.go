package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	version = "1.0.0"
	name    = "gast"
)

var (
	showVersion = flag.Bool("version", false, "显示版本信息")
	showHelp    = flag.Bool("help", false, "显示帮助信息")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if *showHelp || len(os.Args) == 1 {
		printHelp()
		return
	}

	// 获取子命令
	subcommand := os.Args[1]
	args := os.Args[2:]
	
	// 路由到对应的命令处理器
	routeCommand(subcommand, args)
}

// 命令路由器
func routeCommand(subcommand string, args []string) {
	// 尝试基本命令
	if handleBasicCommands(subcommand) {
		return
	}
	
	// 尝试配置命令
	if subcommand == "config" {
		handleConfigCommand(args)
		return
	}
	
	// 尝试文件相关命令
	if handleFileCommands(subcommand, args) {
		return
	}
	
	// 尝试网络相关命令
	if handleNetworkCommands(subcommand, args) {
		return
	}
	
	// 尝试grep命令
	if handleGrepCommands(subcommand, args) {
		return
	}
	
	// 交互模式
	if subcommand == "interactive" {
		handleInteractiveCommand()
		return
	}
	
	// 未知命令
	fmt.Fprintf(os.Stderr, "未知命令: %s\n", subcommand)
	printHelp()
	os.Exit(1)
} 