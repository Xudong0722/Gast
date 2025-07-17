package main

import (
	"fmt"
	"strings"
)

// 哈希命令处理函数
func handleHashCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("用法: gast hash <文件> <类型:md5|sha256>")
		return
	}
	
	filename := args[0]
	hashType := args[1]
	
	hash, err := calculateFileHash(filename, hashType)
	if err != nil {
		fmt.Printf("计算哈希失败: %v\n", err)
		return
	}
	
	fmt.Printf("文件: %s\n", filename)
	fmt.Printf("%s: %s\n", strings.ToUpper(hashType), hash)
}

// 文件查找命令处理函数
func handleFindCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("用法: gast find <目录> <模式>")
		return
	}
	
	findFiles(args[0], args[1])
}

// 目录分析命令处理函数
func handleAnalyzeCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("用法: gast analyze <目录>")
		return
	}
	
	analyzeDirectory(args[0])
}

// 文件处理命令处理函数
func handleProcessCommand(args []string) {
	workers := 4
	dir := "."
	
	if len(args) >= 1 {
		dir = args[0]
	}
	
	if len(args) >= 2 {
		fmt.Sscanf(args[1], "%d", &workers)
	}
	
	processFiles(dir, workers)
}

// 处理文件相关命令
func handleFileCommands(subcommand string, args []string) bool {
	switch subcommand {
	case "hash":
		handleHashCommand(args)
		return true
	case "find":
		handleFindCommand(args)
		return true
	case "analyze":
		handleAnalyzeCommand(args)
		return true
	case "process":
		handleProcessCommand(args)
		return true
	default:
		return false
	}
} 