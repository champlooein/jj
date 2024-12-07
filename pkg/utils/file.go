package utils

import (
	"os"

	"github.com/pkg/errors"
)

func WriteToFile(path string, data string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "file create failed")
	}
	defer f.Close() // 确保在函数结束时关闭文件

	_, err = f.WriteString(data)
	if err != nil {
		return errors.Wrap(err, "write failed: %v")
	}

	return nil
}
