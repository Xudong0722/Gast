package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 交互式输入
func interactiveMode() {
	fmt.Println("进入交互模式 (输入 'quit' 退出)")
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("gast> ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" || input == "exit" {
			fmt.Println("再见!")
			break
		}
		
		if input == "" {
			continue
		}
		
		// 处理交互命令
		parts := strings.Fields(input)
		subcommand := parts[0]
		args := parts[1:]
		
		// 尝试各种命令处理器
		if handleBasicCommands(subcommand) {
			continue
		}
		
		if subcommand == "config" {
			handleConfigCommand(args)
			continue
		}
		
		if handleFileCommands(subcommand, args) {
			continue
		}
		
		if handleNetworkCommands(subcommand, args) {
			continue
		}
		
		if handleGrepCommands(subcommand, args) {
			continue
		}
		
		// 特殊处理grep命令在交互模式中的简化版本
		if subcommand == "grep" {
			if len(parts) < 2 {
				fmt.Println("用法: grep <模式> [文件]")
			} else {
				pattern := parts[1]
				targets := []string{"."}
				if len(parts) > 2 {
					targets = parts[2:]
				}
				options := &GrepOptions{
					IgnoreCase:  false,
					ShowLineNum: true,
					Recursive:   true,
					InvertMatch: false,
					CountOnly:   false,
					FilesOnly:   false,
					Color:       "auto",
					Text:        false,
					Context:     0,
				}
				grepSearch(pattern, targets, options)
			}
			continue
		}
		
		// 未知命令
		fmt.Printf("未知命令: %s\n", input)
		fmt.Println("可用命令: info, version, config, benchmark, hash, url, find, analyze, process, cat, grep, quit")
	}
}

// 处理交互模式命令
func handleInteractiveCommand() {
	interactiveMode()
} 