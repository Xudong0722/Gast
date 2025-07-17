package main

import (
	"fmt"
)

// URL测试命令处理函数
func handleURLCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("用法: gast url <URL>")
		return
	}
	
	testURL(args[0])
}

// 处理网络相关命令
func handleNetworkCommands(subcommand string, args []string) bool {
	switch subcommand {
	case "url":
		handleURLCommand(args)
		return true
	default:
		return false
	}
} 