package main

import (
	"fmt"
	"strconv"
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
		fmt.Println("  -C, --context=NUM    显示匹配行前后各NUM行")
		fmt.Println("  --color[=WHEN]       高亮匹配文本 (auto, always, never)")
		fmt.Println("  --text               强制将二进制文件作为文本处理")
		fmt.Println("示例:")
		fmt.Println("  gast grep -i \"hello\" .")
		fmt.Println("  gast grep -n \"func main\" main.go")
		fmt.Println("  gast grep -r \"TODO\" src/")
		fmt.Println("  gast grep -C 3 \"error\" file.txt")
		fmt.Println("  gast grep --color=auto \"pattern\" file.txt")
		return
	}
	
	options := &GrepOptions{
		IgnoreCase:   false,
		ShowLineNum:  false,
		Recursive:    false,
		InvertMatch:  false,
		CountOnly:    false,
		FilesOnly:    false,
		Color:        "auto",
		Text:         false,
		Context:      0,
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
		case "--color":
			options.Color = "always"
		case "--text":
			options.Text = true
		case "-C":
			// -C 后面应该跟一个数字
			if i+1 >= len(args) {
				fmt.Println("错误: -C 选项需要指定上下文行数")
				return
			}
			i++
			contextStr := args[i]
			contextNum, err := strconv.Atoi(contextStr)
			if err != nil || contextNum < 0 {
				fmt.Printf("错误: 无效的上下文行数: %s\n", contextStr)
				return
			}
			options.Context = contextNum
		default:
			// 检查是否是--color=value格式
			if strings.HasPrefix(arg, "--color=") {
				colorValue := arg[8:] // 去掉"--color="
				if colorValue == "auto" || colorValue == "always" || colorValue == "never" {
					options.Color = colorValue
				} else {
					fmt.Printf("无效的颜色选项: %s (可用: auto, always, never)\n", colorValue)
					return
				}
			} else if strings.HasPrefix(arg, "--context=") {
				contextStr := arg[10:] // 去掉"--context="
				contextNum, err := strconv.Atoi(contextStr)
				if err != nil || contextNum < 0 {
					fmt.Printf("错误: 无效的上下文行数: %s\n", contextStr)
					return
				}
				options.Context = contextNum
			} else {
				fmt.Printf("未知选项: %s\n", arg)
				return
			}
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