package main

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

var (
	Config *config
)

type config struct {
	// 服务端口
	Port string `json:"port"`
	// 服务地址
	Host string `json:"host"`
	// 服务名称
	Name string `json:"name"`
	// 服务版本
	Version    string `json:"version"`
	Jwt_secret string `json:"jwt_secret"`
	// 数据库
	Database struct {
		// 用户名
		Username string `json:"username"`
		// 密码
		Password string `json:"password"`
		// 数据库名
		Database string `json:"database"`
		// 地址
		Host string `json:"host"`
		// 端口
		Port string `json:"port"`
	}
}

func init() {
	// 检查config目录是否存在
	err := os.MkdirAll("./config", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// 检查config.json是否存在
	if _, err := os.Stat("./config/config.json"); os.IsNotExist(err) {
		// 不存在则创建
		file, err := os.Create("./config/config.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		// 写入默认配置
		var defaultConfig = config{}
		defaultConfig.Port = "8080"
		defaultConfig.Host = "127.0.0.1"
		defaultConfig.Name = "My Documents"
		defaultConfig.Version = "1.0.0"
		defaultConfig.Database.Username = "root"
		defaultConfig.Database.Password = "123456"
		defaultConfig.Database.Database = "my_documents"
		defaultConfig.Database.Host = "127.0.0.1"
		defaultConfig.Database.Port = "3306"
		defaultConfig.Jwt_secret = uuid.New().String()

		// 格式化json
		data, err := json.MarshalIndent(defaultConfig, "", "    ")

		// 写入文件
		_, err = file.Write(data)
		if err != nil {
			panic(err)
		}

		// 赋值
		Config = &defaultConfig

		// 关闭文件
		err = file.Close()
		if err != nil {
			return
		}

	} else {
		// 存在则读取
		file, err := os.Open("./config/config.json")
		if err != nil {
			panic(err)
		}

		// 读取文件
		data := make([]byte, 1024)
		n, err := file.Read(data)
		if err != nil {
			panic(err)
		}

		// 解析json
		var config = config{}
		err = json.Unmarshal(data[:n], &config)
		if err != nil {
			panic(err)
		}

		// 如果jwt_secret为空则生成一个
		if config.Jwt_secret == "" {
			config.Jwt_secret = uuid.New().String()
		}

		// 赋值
		Config = &config

		// 关闭文件
		err = file.Close()
		if err != nil {
			return
		}
	}
}

//
