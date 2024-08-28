package main

import (
	"fmt"
	"os"
)

func ensureDirExists(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Println("目录创建成功:", dirPath)
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return err
	} else {
		fmt.Println("目录已存在:", dirPath)
	}

	return nil
}

// 检测文件是否存在
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
