package main

import (
	"fmt"
	"mcsr/src/cmd"
)

func main() {

	// 执行根命令，这会解析命令行输入并运行相应逻辑
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Command execution failed: %v\n", err)
	}

}
