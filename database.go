package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DB 数据库连接
	DB *gorm.DB
)

// 初始化数据库
func InitDB() {
	// 读取配置文件
	username := Config.Database.Username
	password := Config.Database.Password
	database := Config.Database.Database
	host := Config.Database.Host
	port := Config.Database.Port

	// 连接数据库
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	fmt.Println("数据库连接成功")

	checkDatabase()
}

func checkDatabase() {
	err := DB.AutoMigrate(
		&User{},
		&Token{},

		&Document{},
		&DocumentVersion{},
		&MarkdownDraft{},
		&Book{},
		&Category{},
		&DocumentList{},
	)
	if err != nil {
		panic(err)
	}
}
