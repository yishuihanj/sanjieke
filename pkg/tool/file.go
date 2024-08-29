package tool

import (
	"io/fs"
	"os"
	"path/filepath"
)

func DeleteFilesInDir(dir string) error {
	// 遍历目录下的所有文件
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 如果不是目录，删除文件
		if !d.IsDir() {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func EnsureDirExists(dirPath string) (error, bool) {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err, false
		}
	} else if err != nil {
		return err, false
	} else if !info.IsDir() {
		return err, false
	} else {
		return nil, true
	}

	return nil, false
}

// 检测文件是否存在
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
