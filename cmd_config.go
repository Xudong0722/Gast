package main

import (
	"fmt"
)

// 配置命令处理函数
func handleConfigCommand(args []string) {
	if len(args) == 0 {
		printConfig()
		return
	}
	
	switch args[0] {
	case "init":
		initConfig()
	case "show":
		printConfig()
	default:
		fmt.Printf("未知配置命令: %s\n", args[0])
		fmt.Println("可用命令: init, show")
	}
} 