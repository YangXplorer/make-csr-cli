package cmd

import (
	"fmt"
	"mcsr/src/config"
	"mcsr/src/internal"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func Execute() error {
	config.InitConfig()
	return rootCmd.Execute()
}

// 创建一个根命令，用于 CLI 工具的入口点
var rootCmd = &cobra.Command{
	Use:   "makeCsr",                                         // CLI 工具的名称
	Short: "A simple CLI tool to generate OpenSSL CSR files", // 工具的简短描述
	Run:   runGenerateCSR,                                    // 当用户输入根命令时，执行 runGenerateCSR 函数
}

func runGenerateCSR(cmd *cobra.Command, args []string) {
	configDir, configFile, cnFile, err := config.InitPaths()
	if err != nil {
		fmt.Printf("failed to get config paths: %v\n", err)
		return
	}

	var cn string
	for {
		// 从文件中读取 Common Name (CN) 选项列表
		commonNames, err := internal.LoadCommonNames(cnFile)
		if err != nil || len(commonNames) == 0 {
			fmt.Println("cn.txt is empty or missing. Please add Common Name (CN) entries to cn.txt and try again.")
			return
		}
		commonNames = append(commonNames, "Enter a new CN.")
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
		_, cn, err = prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		// 如果用户选择了 "Add new CN" 选项，提示用户输入新的 CN
		if cn == "Enter a new CN." {
			if _, err = config.EnsureFileWithContent(false); err != nil {
				fmt.Printf("failed to ensure cn.txt file: %v\n", err)
				continue
			}
		} else {
			// 用户选择了其他 CN，退出循环
			fmt.Printf("You selected: %s\n", cn)
			break
		}
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
}
