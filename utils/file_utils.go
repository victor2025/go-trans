package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func LoadJsonFile(path string, v *map[string]any) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析 JSON
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
