package util

import (
	"crypto/md5"
	"fmt"
	"github.com/cilidm/toolbox/logging"
	"io"
	"os"
)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetMdByPath(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		logging.Error(err.Error())
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		logging.Error(err.Error())
		return "", err
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}
