package internal

import (
	"bufio"
	"os"
	"strings"
)

// LoadCommonNames 从指定的文件中读取 Common Name (CN) 列表
func LoadCommonNames(filename string) ([]string, error) {
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
