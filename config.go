package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	LogLevel    string `json:"log_level"`
	OutputDir   string `json:"output_dir"`
	MaxWorkers  int    `json:"max_workers"`
	Timeout     int    `json:"timeout"`
	EnableColor bool   `json:"enable_color"`
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".gast.json"
	}
	return filepath.Join(home, ".gast.json")
}

func loadConfig() (*Config, error) {
	config := &Config{
		LogLevel:    "info",
		OutputDir:   "./output",
		MaxWorkers:  4,
		Timeout:     30,
		EnableColor: true,
	}

	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return config, nil
}

func saveConfig(config *Config) error {
	configPath := getConfigPath()
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	return nil
}

func printConfig() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	fmt.Println("当前配置:")
	fmt.Printf("  日志级别: %s\n", config.LogLevel)
	fmt.Printf("  输出目录: %s\n", config.OutputDir)
	fmt.Printf("  最大工作线程: %d\n", config.MaxWorkers)
	fmt.Printf("  超时时间: %d 秒\n", config.Timeout)
	fmt.Printf("  启用颜色: %v\n", config.EnableColor)
	fmt.Printf("  配置文件: %s\n", getConfigPath())
}

func initConfig() {
	config := &Config{
		LogLevel:    "info",
		OutputDir:   "./output",
		MaxWorkers:  4,
		Timeout:     30,
		EnableColor: true,
	}

	if err := saveConfig(config); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		return
	}

	fmt.Printf("配置文件已创建: %s\n", getConfigPath())
	printConfig()
} 