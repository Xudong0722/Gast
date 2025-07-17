package main

import (
	"fmt"
	"strings"
)

// Grep命令处理函数
func handleGrepCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("用法: gast grep [选项] <模式> [文件/目录]")
		fmt.Println("选项:")
		fmt.Println("  -i, --ignore-case    忽略大小写")
		fmt.Println("  -n, --line-number    显示行号")
		fmt.Println("  -r, --recursive      递归搜索目录")
		fmt.Println("  -v, --invert-match   反向匹配")
		fmt.Println("  -c, --count          只显示匹配行数")
		fmt.Println("  -l, --files-with-matches  只显示匹配的文件名")
		fmt.Println("示例:")
		fmt.Println("  gast grep -i \"hello\" .")
		fmt.Println("  gast grep -n \"func main\" main.go")
		fmt.Println("  gast grep -r \"TODO\" src/")
		return
	}
	
	options := &GrepOptions{
		IgnoreCase:   false,
		ShowLineNum:  false,
		Recursive:    false,
		InvertMatch:  false,
		CountOnly:    false,
		FilesOnly:    false,
	}
	
	var pattern string
	var targets []string
	
	// 解析参数
	i := 0
	for i < len(args) {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			pattern = arg
			i++
			break
		}
		
		switch arg {
		case "-i", "--ignore-case":
			options.IgnoreCase = true
		case "-n", "--line-number":
			options.ShowLineNum = true
		case "-r", "--recursive":
			options.Recursive = true
		case "-v", "--invert-match":
			options.InvertMatch = true
		case "-c", "--count":
			options.CountOnly = true
		case "-l", "--files-with-matches":
			options.FilesOnly = true
		default:
			fmt.Printf("未知选项: %s\n", arg)
			return
		}
		i++
	}
	
	// 收集目标文件/目录
	for i < len(args) {
		targets = append(targets, args[i])
		i++
	}
	
	// 如果没有指定目标，默认为当前目录
	if len(targets) == 0 {
		targets = []string{"."}
	}
	
	if pattern == "" {
		fmt.Println("错误: 必须指定搜索模式")
		return
	}
	
	grepSearch(pattern, targets, options)
}

// 处理grep相关命令
func handleGrepCommands(subcommand string, args []string) bool {
	switch subcommand {
	case "grep":
		handleGrepCommand(args)
		return true
	default:
		return false
	}
} 