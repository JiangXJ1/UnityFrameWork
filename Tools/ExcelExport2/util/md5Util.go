package util

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// 读取文件并将内容存入map
func LoadHashesFromFile(filePath string) (map[string]string, error) {
	hashMap := make(map[string]string)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "===")
		if len(parts) == 2 {
			hashMap[parts[0]] = parts[1]
		} else {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件MD5失败")
		return hashMap, err
	}

	return hashMap, nil
}

// WriteHashesToFile 将map内容写回到文件
func WriteHashesToFile(filePath string, hashMap map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭文件失败")
		}
	}(file)

	writer := bufio.NewWriter(file)
	for name, hash := range hashMap {
		_, err := writer.WriteString(fmt.Sprintf("%s===%s\n", name, hash))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// 计算文件MD5
func CalculateFileHash(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	_, err = hasher.Write(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
