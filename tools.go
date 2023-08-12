package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
)

// saveFile 保存文件
func saveFile(filePath string, content string) error {

	// 创建/files目录
	err := os.MkdirAll("./files", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// 写入文件
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	// 关闭文件
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// sha256File 计算文件的sha256值
func sha256File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum), nil
}

// GetFileType 获取文件类型
func GetFileType(fileName string) string {
	// 获取文件后缀
	fileSuffix := path.Ext(fileName)
	switch fileSuffix {
	case ".md":
		return "markdown"
	case ".docx":
		return "docx"
	case ".xlsx":
		return "xlsx"
	case ".pptx":
		return "pptx"
	case ".pdf":
		return "pdf"
	default:
		return "unknown"
	}
}
