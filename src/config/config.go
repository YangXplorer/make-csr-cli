package config

import (
	"fmt"
	"mcsr/src/internal"
	"os"
	"path/filepath"
	"runtime"
)

// InitPaths 初始化配置文件的相关路径，返回配置目录、配置文件路径和 CN 文件路径
func InitPaths() (configDir, configFile, cnFile string, err error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// 根据操作系统设置路径
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(homeDir, ".config", "openssl")
	} else {
		configDir = filepath.Join(homeDir, ".config/openssl")
	}
	configFile = filepath.Join(configDir, "openssl.conf")
	cnFile = filepath.Join(configDir, "cn.txt")

	return configDir, configFile, cnFile, nil
}

// InitConfig 确保配置文件目录、config 文件和 cn.txt 文件存在
func InitConfig() error {
	configDir, configFile, cnFile, err := InitPaths()
	if err != nil {
		return fmt.Errorf("failed to get config paths: %v", err)
	}

	// 确保配置目录存在
	if err := ensureDirectory(configDir); err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// 确保配置文件存在
	if err := ensureConfigFile(configFile); err != nil {
		return fmt.Errorf("failed to ensure config file: %w", err)
	}

	// 确保 CN 文件存在
	if err := ensureFileWithContent(cnFile, "www.example.com\n"); err != nil {
		return fmt.Errorf("failed to ensure cn.txt file: %w", err)
	}

	return nil
}

// ensureDirectory 确保目录存在，如果不存在则创建
func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", path, err)
		}
	}
	return nil
}

// ensureConfigFile 确保配置文件存在，如果不存在则引导用户创建
func ensureConfigFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Configuration file not found. Please provide the following details:")

		// 提示用户输入各个字段
		country := internal.UserCommandInput("Country Code (C) [e.g., JP]: ", "JP")
		state := internal.UserCommandInput("State or Province Name (ST) [e.g., Tokyo]: ", "Tokyo")
		locality := internal.UserCommandInput("Locality Name (L) [e.g., CHUOU-KU]: ", "CHUOU-KU")
		organization := internal.UserCommandInput("Organization Name (O) [e.g., BRIDGE CO.,LTD.]: ", "BRIDGE CO.,LTD.")
		orgUnit := internal.UserCommandInput("Organizational Unit Name (OU) [e.g., BRIDGE CO.,LTD.]: ", "BRIDGE CO.,LTD.")

		// 构建配置文件内容
		configContent := fmt.Sprintf(`
[ req ]
default_bits        = 2048
default_keyfile     = openssl-default.key
distinguished_name  = req_distinguished_name
prompt              = no

[ req_distinguished_name ]
C  = %s
ST = %s
L  = %s
O  = %s
OU = %s
`, country, state, locality, organization, orgUnit)

		// 写入配置文件
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create config file '%s': %w", path, err)
		}
		defer file.Close()

		if _, err := file.WriteString(configContent); err != nil {
			return fmt.Errorf("failed to write to config file '%s': %w", path, err)
		}
	}
	return nil
}

// ensureFileWithContent 确保文件存在，如果不存在则创建并写入默认内容
func ensureFileWithContent(path, defaultContent string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file '%s': %w", path, err)
		}
		defer file.Close()

		if _, err := file.WriteString(defaultContent); err != nil {
			return fmt.Errorf("failed to write default content to file '%s': %w", path, err)
		}
	}
	return nil
}
