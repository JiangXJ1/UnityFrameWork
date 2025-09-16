package fileUtil

import (
	"bufio"
	"crypto/md5"
	"demo-go-excel/global"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DeleteDir(dirPath string) {
	os.RemoveAll(dirPath)
}

// 按过滤列表删除文件
func DeleteFiles(dirPath string, filterList map[string]bool) {
	files, _ := os.ReadDir(dirPath)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileName := f.Name()
		if filterList == nil || len(filterList) == 0 {
			DeleteFile(dirPath + fileName)
			continue
		}
		if filterList[fileName] {
			DeleteFile(dirPath + fileName)
		}
	}
}

func DeleteFile(fileName string) {
	var err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func CreateDir(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func FileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateFile(filePath string) {
	if FileExist(filePath) {
		DeleteFile(filePath)
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

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

// 将map内容写回到文件
func WriteHashesToFile(filePath string, hashMap map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for name, hash := range hashMap {
		_, err := writer.WriteString(fmt.Sprintf("%s===%s\n", name, hash))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func GetCSFilePathName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "@", "")
	return filepath.Join(global.OUTPUT_DIR_CS, "Schema_"+fileName+".cs")
}

func GetBinFilePathName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "@", "")
	return filepath.Join(global.OUTPUT_DIR_BIN, "Schema_"+fileName+".bytes")
}

func GetServerLuaFilePathName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "@", "")
	return filepath.Join(global.OUTPUT_DIR_LUA, "Schema_"+fileName+".lua")
}

func GetServerCSVFilePathName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "@", "")
	return filepath.Join(global.OUTPUT_DIR_CSV, "Schema_"+fileName+".csv")
}

func CreateDirectoryIfNotExists(path string) {
	_, err := os.Stat(path)
	if err == nil {
		os.RemoveAll(path)
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("无法创建目录: %s", err)
	}

}
