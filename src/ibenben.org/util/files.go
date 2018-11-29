package util

import (
	"os"
	"io"
)


//文件拷贝
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}


//创建文件或者目录（如果不存在）
func CreateFile(name string) error {
	if !IsFileExists(name) {
		err := os.MkdirAll(name, os.ModePerm)
		return err
	}
	return nil
}

//文件或者目录是否存在
func IsFileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
