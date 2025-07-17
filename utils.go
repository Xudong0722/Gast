package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ANSI颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// 检查终端是否支持颜色
func isTerminalColorSupported() bool {
	// 检查是否是终端
	if !isTerminal(os.Stdout) {
		return false
	}
	
	// Windows平台特殊处理
	if runtime.GOOS == "windows" {
		// 检查是否是Windows Terminal或新版本的CMD/PowerShell
		colorterm := os.Getenv("COLORTERM")
		if colorterm != "" {
			return true
		}
		
		// 检查WT_SESSION环境变量（Windows Terminal）
		if os.Getenv("WT_SESSION") != "" {
			return true
		}
		
		// 检查ConEmu环境变量
		if os.Getenv("ConEmuANSI") == "ON" {
			return true
		}
		
		// 对于传统的Windows CMD/PowerShell，检查是否支持ANSI
		// 如果没有特殊环境变量，我们假设支持（现代Windows版本都支持）
		return true
	}
	
	// 非Windows平台，检查TERM环境变量
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	
	return true
}

// 检查文件描述符是否是终端
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	
	// Windows平台特殊处理
	if runtime.GOOS == "windows" {
		// 在Windows上，检查是否是字符设备
		mode := fi.Mode()
		return mode&os.ModeCharDevice != 0
	}
	
	// 非Windows平台，使用标准检查
	return fi.Mode()&os.ModeCharDevice != 0
}

// 应该使用颜色
func shouldUseColor(colorOption string) bool {
	switch colorOption {
	case "always":
		return true
	case "never":
		return false
	case "auto":
		return isTerminalColorSupported()
	default:
		return false
	}
}

// Grep选项
type GrepOptions struct {
	IgnoreCase   bool
	ShowLineNum  bool
	Recursive    bool
	InvertMatch  bool
	CountOnly    bool
	FilesOnly    bool
	Color        string // "auto", "always", "never"
	Text         bool   // 强制将二进制文件作为文本处理
}

// Grep搜索结果
type GrepResult struct {
	Filename string
	LineNum  int
	Line     string
	Matches  []string
}

// 文件哈希计算
func calculateFileHash(filename string, hashType string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	switch hashType {
	case "md5":
		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return "", err
		}
		return fmt.Sprintf("%x", hash.Sum(nil)), nil
	case "sha256":
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return "", err
		}
		return fmt.Sprintf("%x", hash.Sum(nil)), nil
	default:
		return "", fmt.Errorf("不支持的哈希类型: %s", hashType)
	}
}

// URL验证和测试
func testURL(target string) {
	// 验证URL格式
	parsedURL, err := url.Parse(target)
	if err != nil {
		fmt.Printf("无效的URL格式: %v\n", err)
		return
	}

	if parsedURL.Scheme == "" {
		target = "http://" + target
	}

	fmt.Printf("测试URL: %s\n", target)
	
	start := time.Now()
	resp, err := http.Get(target)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ 响应时间: %v\n", duration)
	fmt.Printf("   状态码: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("   内容长度: %d bytes\n", resp.ContentLength)
	fmt.Printf("   内容类型: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("   服务器: %s\n", resp.Header.Get("Server"))
}

// 文件查找
func findFiles(dir string, pattern string) {
	fmt.Printf("在 %s 中查找匹配 '%s' 的文件:\n", dir, pattern)
	
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), strings.ToLower(pattern)) {
			fmt.Printf("  %s (%d bytes)\n", path, info.Size())
			count++
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("查找失败: %v\n", err)
		return
	}
	
	fmt.Printf("共找到 %d 个文件\n", count)
}

// 文件大小统计
func analyzeDirectory(dir string) {
	fmt.Printf("分析目录: %s\n", dir)
	
	var totalSize int64
	var fileCount int
	var dirCount int
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			dirCount++
		} else {
			fileCount++
			totalSize += info.Size()
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("分析失败: %v\n", err)
		return
	}
	
	fmt.Printf("结果:\n")
	fmt.Printf("  文件数: %d\n", fileCount)
	fmt.Printf("  目录数: %d\n", dirCount)
	fmt.Printf("  总大小: %.2f MB\n", float64(totalSize)/1024/1024)
}

// 并发文件处理示例
func processFiles(dir string, workers int) {
	fmt.Printf("使用 %d 个工作线程处理文件...\n", workers)
	
	filesChan := make(chan string, 100)
	resultsChan := make(chan string, 100)
	var wg sync.WaitGroup
	
	// 启动工作线程
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for filename := range filesChan {
				// 模拟文件处理
				time.Sleep(100 * time.Millisecond)
				resultsChan <- fmt.Sprintf("工作线程 %d 处理: %s", id, filename)
			}
		}(i)
	}
	
	// 发送文件到工作线程
	go func() {
		defer close(filesChan)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filesChan <- path
			}
			return nil
		})
	}()
	
	// 收集结果
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	count := 0
	for result := range resultsChan {
		fmt.Printf("  %s\n", result)
		count++
	}
	
	fmt.Printf("处理完成，共处理 %d 个文件\n", count)
}



// Grep搜索主函数
func grepSearch(pattern string, targets []string, options *GrepOptions) {
	// 编译正则表达式
	var regex *regexp.Regexp
	var err error
	
	if options.IgnoreCase {
		regex, err = regexp.Compile("(?i)" + pattern)
	} else {
		regex, err = regexp.Compile(pattern)
	}
	
	if err != nil {
		fmt.Printf("正则表达式编译错误: %v\n", err)
		return
	}
	
	totalMatches := 0
	
	for _, target := range targets {
		matches := processGrepTarget(target, regex, options)
		totalMatches += matches
	}
	
	if options.CountOnly {
		fmt.Printf("总匹配数: %d\n", totalMatches)
	}
}

// 处理grep目标（文件或目录）
func processGrepTarget(target string, regex *regexp.Regexp, options *GrepOptions) int {
	info, err := os.Stat(target)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return 0
	}
	
	if info.IsDir() {
		if options.Recursive {
			return grepInDirectory(target, regex, options)
		} else {
			fmt.Printf("跳过目录: %s (使用 -r 选项递归搜索)\n", target)
			return 0
		}
	} else {
		if options.Text || isTextFile(target) {
			return grepInFile(target, regex, options)
		} else {
			return 0
		}
	}
}

// 在目录中递归搜索
func grepInDirectory(dir string, regex *regexp.Regexp, options *GrepOptions) int {
	totalMatches := 0
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && (options.Text || isTextFile(path)) {
			matches := grepInFile(path, regex, options)
			totalMatches += matches
		}
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("遍历目录错误: %v\n", err)
	}
	
	return totalMatches
}

// 在单个文件中搜索
func grepInFile(filename string, regex *regexp.Regexp, options *GrepOptions) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("打开文件错误 %s: %v\n", filename, err)
		return 0
	}
	defer file.Close()
	
	reader := bufio.NewReader(file)
	lineNum := 0
	matchCount := 0
	hasMatch := false
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 处理文件末尾没有换行符的情况
				if len(line) > 0 {
					lineNum++
					processLine(line, filename, lineNum, regex, options, &matchCount, &hasMatch)
				}
				break
			}
			fmt.Printf("读取文件错误 %s: %v\n", filename, err)
			break
		}
		
		lineNum++
		// 移除行尾的换行符
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		
		processLine(line, filename, lineNum, regex, options, &matchCount, &hasMatch)
	}
	
	// 输出统计信息
	if options.CountOnly {
		fmt.Printf("%s: %d\n", filename, matchCount)
	} else if options.FilesOnly && hasMatch {
		fmt.Println(filename)
	}
	
	return matchCount
}

// 处理单行匹配
func processLine(line string, filename string, lineNum int, regex *regexp.Regexp, options *GrepOptions, matchCount *int, hasMatch *bool) {
	matches := regex.FindAllString(line, -1)
	isMatch := len(matches) > 0
	
	// 处理反向匹配
	if options.InvertMatch {
		isMatch = !isMatch
	}
	
	if isMatch {
		*matchCount++
		*hasMatch = true
		
		if !options.CountOnly && !options.FilesOnly {
			result := &GrepResult{
				Filename: filename,
				LineNum:  lineNum,
				Line:     line,
				Matches:  matches,
			}
			printGrepResult(result, options)
		}
	}
}

// 打印grep结果
func printGrepResult(result *GrepResult, options *GrepOptions) {
	var output strings.Builder
	useColor := shouldUseColor(options.Color)
	
	// 文件名 (紫色)
	if useColor {
		output.WriteString(ColorPurple + result.Filename + ColorReset)
	} else {
		output.WriteString(result.Filename)
	}
	
	// 行号 (绿色)
	if options.ShowLineNum {
		if useColor {
			output.WriteString(":" + ColorGreen + fmt.Sprintf("%d", result.LineNum) + ColorReset)
		} else {
			output.WriteString(fmt.Sprintf(":%d", result.LineNum))
		}
	}
	
	output.WriteString(": ")
	
	// 高亮匹配的文本
	line := result.Line
	if len(result.Matches) > 0 && !options.InvertMatch {
		// 为了避免重复替换，我们需要按长度排序，从长到短进行替换
		matches := make([]string, len(result.Matches))
		copy(matches, result.Matches)
		
		// 简单的按长度排序（冒泡排序）
		for i := 0; i < len(matches); i++ {
			for j := i + 1; j < len(matches); j++ {
				if len(matches[i]) < len(matches[j]) {
					matches[i], matches[j] = matches[j], matches[i]
				}
			}
		}
		
		if useColor {
			// 使用ANSI颜色代码高亮匹配的文本
			for _, match := range matches {
				highlightedMatch := ColorRed + ColorBold + match + ColorReset
				line = strings.ReplaceAll(line, match, highlightedMatch)
			}
		} else {
			// 使用方括号包围匹配的文本
			for _, match := range matches {
				line = strings.ReplaceAll(line, match, "["+match+"]")
			}
		}
	}
	
	output.WriteString(line)
	fmt.Println(output.String())
}

// 判断是否为文本文件
func isTextFile(filename string) bool {
	// 基于文件扩展名的简单判断
	ext := strings.ToLower(filepath.Ext(filename))
	textExtensions := []string{
		".txt", ".md", ".go", ".py", ".java", ".c", ".cpp", ".h", ".hpp",
		".js", ".ts", ".html", ".css", ".json", ".xml", ".yaml", ".yml",
		".sh", ".bash", ".zsh", ".fish", ".ps1", ".bat", ".cmd",
		".sql", ".php", ".rb", ".rs", ".swift", ".kt", ".scala",
		".r", ".m", ".pl", ".lua", ".vim", ".emacs", ".cfg", ".conf",
		".ini", ".toml", ".properties", ".log", ".csv", ".tsv",
	}
	
	for _, validExt := range textExtensions {
		if ext == validExt {
			return true
		}
	}
	
	// 如果没有扩展名，检查文件内容
	if ext == "" {
		return isTextContent(filename)
	}
	
	return false
}

// 检查文件内容是否为文本
func isTextContent(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()
	
	// 读取前512字节来判断
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}
	
	// 检查是否包含二进制字符
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return false
		}
		if buffer[i] < 32 && buffer[i] != 9 && buffer[i] != 10 && buffer[i] != 13 {
			return false
		}
	}
	
	return true
} 