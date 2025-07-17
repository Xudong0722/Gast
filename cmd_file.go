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

// Cat命令处理函数
func handleCatCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("用法: gast cat [选项] <文件1> [文件2] ...")
		fmt.Println("选项:")
		fmt.Println("  -n, --number             显示行号")
		fmt.Println("  -b, --number-nonblank    显示非空行的行号")
		fmt.Println("  -E, --show-ends          在行尾显示$符号")
		fmt.Println("  -T, --show-tabs          显示制表符为^I")
		fmt.Println("  -A, --show-all           显示所有字符 (相当于-vET)")
		fmt.Println("  -v, --show-nonprinting   显示非打印字符")
		fmt.Println("示例:")
		fmt.Println("  gast cat file.txt")
		fmt.Println("  gast cat -n file1.txt file2.txt")
		fmt.Println("  gast cat -A file.txt")
		return
	}
	
	options := &CatOptions{
		ShowLineNumbers:   false,
		ShowNonEmpty:      false,
		ShowEnds:          false,
		ShowTabs:          false,
		ShowAll:           false,
		ShowNonPrinting:   false,
	}
	
	var filenames []string
	
	// 解析参数
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			filenames = args[i:]
			break
		}
		
		switch arg {
		case "-n", "--number":
			options.ShowLineNumbers = true
		case "-b", "--number-nonblank":
			options.ShowNonEmpty = true
		case "-E", "--show-ends":
			options.ShowEnds = true
		case "-T", "--show-tabs":
			options.ShowTabs = true
		case "-A", "--show-all":
			options.ShowAll = true
			options.ShowNonPrinting = true
			options.ShowEnds = true
			options.ShowTabs = true
		case "-v", "--show-nonprinting":
			options.ShowNonPrinting = true
		default:
			fmt.Printf("未知选项: %s\n", arg)
			return
		}
	}
	
	if len(filenames) == 0 {
		fmt.Println("错误: 必须指定至少一个文件")
		return
	}
	
	// 处理每个文件
	for _, filename := range filenames {
		if len(filenames) > 1 {
			fmt.Printf("==> %s <==\n", filename)
		}
		
		if err := catFile(filename, options); err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}
		
		if len(filenames) > 1 {
			fmt.Println()
		}
	}
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
	case "cat":
		handleCatCommand(args)
		return true
	default:
		return false
	}
} 