package util

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFolder 复制文件夹及其内容
func CopyFolder(src string, dst string) error {
	// 创建目标文件夹
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	// 读取源文件夹中的所有文件和子文件夹
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 遍历源文件夹的每个条目
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// 如果是目录，则递归调用 CopyFolder
		if entry.IsDir() {
			if err := CopyFolder(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 如果是文件，则复制文件
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile 复制单个文件
func CopyFile(src string, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 复制文件内容
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
