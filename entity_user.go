package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(255);not null;comment:用户名"`
	Password string `gorm:"type:varchar(255);not null;comment:密码"`
	Email    string `gorm:"type:varchar(255);not null;comment:邮箱"`
	Avatar   string `gorm:"type:varchar(255);not null;comment:头像"`
	Role     int    `gorm:"type:int;not null;comment:角色"`
}

type Token struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	UserID uint   `gorm:"type:int;not null;comment:用户id"`
	Token  string `gorm:"type:varchar(255);not null;comment:token"`
	Expire int64  `gorm:"type:bigint;not null;comment:过期时间"`
	Status int    `gorm:"type:int;not null;comment:状态;default:1"`
}
