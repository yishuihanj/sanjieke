package tool

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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

// 判断目录是否存在，如果不存在则创建
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
func CheckFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// makeValidFilename 将不合法的文件名字符替换为合法字符（下划线）
func MakeValidFilename(filename string) string {
	var result strings.Builder

	for _, r := range filename {
		if isValidFilenameChar(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}

	return result.String()
}

// isValidFilenameChar 判断字符是否为合法的文件名字符
func isValidFilenameChar(r rune) bool {
	// 合法字符可以包括字母、数字和一些符号（例如 -、_ 和 .）
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.'
}
