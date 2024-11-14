package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// UserCommandInput 提示用户输入，如果用户按回车则返回默认值
func UserCommandInput(prompt string, defaultValue string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}
