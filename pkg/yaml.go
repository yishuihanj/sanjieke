package pkg

import (
	"gopkg.in/yaml.v3"
	"os"
)

func YamlReader(configPath string, c interface{}) error {
	// 读取 YAML 文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	// 解析 YAML 数据
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}
