package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func main() {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get user home directory: %v\n", err)
		return
	}

	// 根据操作系统定义配置文件目录和路径
	var configDir, configFile, cnFile string
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(homeDir, ".config\\openssl\\")
		configFile = filepath.Join(configDir, "openssl.conf")
		cnFile = filepath.Join(configDir, "cn.txt")
	} else {
		configDir = filepath.Join(homeDir, ".config/openssl/")
		configFile = filepath.Join(configDir, "openssl.conf")
		cnFile = filepath.Join(configDir, "cn.txt")
	}

	// 确保配置文件和 cn.txt 文件存在
	err = ensureConfigFile(configDir, configFile, cnFile)
	if err != nil {
		fmt.Printf("Failed to create config file or cn.txt: %v\n", err)
		return
	}

	// 创建一个根命令，用于 CLI 工具的入口点
	var rootCmd = &cobra.Command{
		Use:   "makeCsr",                                         // CLI 工具的名称
		Short: "A simple CLI tool to generate OpenSSL CSR files", // 工具的简短描述
		Run: func(cmd *cobra.Command, args []string) {

			// 从文件中读取 Common Name (CN) 选项列表
			commonNames, err := readCommonNamesFromFile(cnFile)
			if err != nil || len(commonNames) == 0 {
				fmt.Println("cn.txt is empty or missing. Please add Common Name (CN) entries to cn.txt and try again.")
				return
			}

			// 自定义 promptui 的选择项样式
			prompt := promptui.Select{
				Label: "コモンネーム (CN) を選択してください。（Select the Common Name (CN)）",
				Items: commonNames,
				Templates: &promptui.SelectTemplates{
					Active:   `{{ "> " | green }}{{ . | green }}`,               // 选中项的前缀和显示样式
					Inactive: `  {{ . }}`,                                       // 非选中项的显示样式
					Selected: `{{ "> " | green | bold }}{{ . | green | bold }}`, // 选中后确认显示样式
				},
				Pointer: promptui.PipeCursor, // 使用 | 作为光标
			}

			// 获取用户选择的索引和值
			_, cn, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			// 提示用户输入 emailAddress，并添加邮件地址格式检查
			validateEmail := func(input string) error {
				// 使用正则表达式检查邮箱格式
				re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !re.MatchString(input) {
					return fmt.Errorf("無効なメールアドレスです（Invalid email address）")
				}
				return nil
			}

			promptEmail := promptui.Prompt{
				Label:    "承認メールアドレスを入れてください。（Enter the email address）",
				Validate: validateEmail,
			}

			email, err := promptEmail.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			// 构建 -subj 参数字符串
			subj := fmt.Sprintf("/C=JP/ST=Tokyo/L=CHUOU-KU/O=BRIDGE CO.,LTD./OU=BRIDGE CO.,LTD./CN=%s/emailAddress=%s", cn, email)

			// 定义 OpenSSL 命令及其参数，用于生成 CSR 和私钥文件
			keyFile := filepath.Join(configDir, cn+".key")
			csrFile := filepath.Join(configDir, cn+".csr")

			// 定义 OpenSSL 命令及其参数，用于生成 CSR 和私钥文件
			opensslCmd := exec.Command("openssl", "req", "-new", "-newkey", "rsa:2048", "-nodes", "-keyout", keyFile, "-out", csrFile, "-subj", subj, "-config", configFile)

			// 设置命令的标准输出和标准错误输出，使其直接显示在终端上
			opensslCmd.Stdout = os.Stdout
			opensslCmd.Stderr = os.Stderr

			// 运行 OpenSSL 命令并检查是否有错误
			if err := opensslCmd.Run(); err != nil {
				// 如果执行命令时出错，打印错误信息
				fmt.Println("Error executing command:", err)
				return
			}

			// 如果成功生成 CSR 文件，打印确认信息
			fmt.Println("CSR has been generated and saved as", cn+".csr")
		},
	}

	// 执行根命令，这会解析命令行输入并运行相应逻辑
	rootCmd.Execute()
}

// ensureConfigFile 确保配置文件目录、config 文件和 cn.txt 文件存在
func ensureConfigFile(configDir, configFile, cnFile string) error {
	// 如果目录不存在，则创建
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// 如果 config 文件不存在，则提示用户输入并创建文件
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("Configuration file not found. Please provide the following details.")

		// 提示用户输入各个字段的值
		country := promptUserInput("Country Code (C) [e.g., JP]: ", "JP")
		state := promptUserInput("State or Province Name (ST) [e.g., Tokyo]: ", "Tokyo")
		locality := promptUserInput("Locality Name (L) [e.g., CHUOU-KU]: ", "CHUOU-KU")
		organization := promptUserInput("Organization Name (O) [e.g., MAKE-CSR CO.,LTD.]: ", "MAKE-CSR CO.,LTD.")
		orgUnit := promptUserInput("Organizational Unit Name (OU) [e.g., MAKE-CSR CO.,LTD.]: ", "MAKE-CSR CO.,LTD.")

		// 根据用户输入构建配置内容
		configContent := fmt.Sprintf(`[ req ]
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

		// 创建并写入配置文件
		file, err := os.Create(configFile)
		if err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(configContent)
		if err != nil {
			return fmt.Errorf("failed to write config file: %v", err)
		}
	}

	// 如果 cn.txt 文件不存在，则创建一个包含默认值的文件
	if _, err := os.Stat(cnFile); os.IsNotExist(err) {
		file, err := os.Create(cnFile)
		if err != nil {
			return fmt.Errorf("failed to create cn.txt file: %v", err)
		}
		defer file.Close()

		// 写入默认的 Common Name
		_, err = file.WriteString("www.example.com\n")
		if err != nil {
			return fmt.Errorf("failed to write to cn.txt file: %v", err)
		}
	}

	return nil
}

// promptUserInput 提示用户输入，如果用户按回车则返回默认值
func promptUserInput(prompt string, defaultValue string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

// readCommonNamesFromFile 从指定的文件中读取 Common Name (CN) 列表
func readCommonNamesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commonNames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			commonNames = append(commonNames, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commonNames, nil
}
