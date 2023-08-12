package main

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"type:varchar(255);not null;comment:文档名称"`
	UUID    string `gorm:"type:varchar(255);not null;comment:文档uuid"`
	Type    string `gorm:"type:varchar(255);not null;comment:文档类型"`
	Path    string `gorm:"type:varchar(255);not null;comment:文档路径"`
	Creator uint   `gorm:"type:int;not null;comment:创建者"`
}

type DocumentVersion struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	DocumentID uint   `gorm:"type:int;not null;comment:文档id"`
	Version    string `gorm:"type:varchar(255);not null;comment:文档版本"`
	UUID       string `gorm:"type:varchar(255);not null;comment:文档版本uuid"`
	Type       string `gorm:"type:varchar(255);not null;comment:文档类型"`
	Size       int64  `gorm:"type:bigint;not null;comment:文档大小"`
	Path       string `gorm:"type:varchar(255);not null;comment:文档路径"`
	Creator    uint   `gorm:"type:int;not null;comment:创建者"`
	SHA256     string `gorm:"type:varchar(255);not null;comment:文档sha256"`
}

// Markdown草稿
type MarkdownDraft struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	UUID    string `gorm:"type:varchar(255);not null;comment:草稿uuid"`
	Content string `gorm:"type:longtext;not null;comment:草稿内容"`
	Creator uint   `gorm:"type:int;not null;comment:创建者"`
}

type Book struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null;comment:书籍名称"`
	Description string `gorm:"type:varchar(255);comment:书籍描述"`
	UUID        string `gorm:"type:varchar(255);not null;comment:书籍uuid"`
	Type        string `gorm:"type:varchar(255);comment:书籍类型"`
	Creator     uint   `gorm:"type:int;not null;comment:创建者"`
}

type Category struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"type:varchar(255);not null;comment:分类名称"`
	UUID    string `gorm:"type:varchar(255);not null;comment:分类uuid"`
	Weight  int    `gorm:"type:int;not null;comment:权重"`
	BookID  uint   `gorm:"type:int;not null;comment:书籍id"`
	Creator uint   `gorm:"type:int;not null;comment:创建者"`
}

type DocumentList struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement"`
	DocumentID uint `gorm:"type:int;not null;comment:文档id"`
	BookID     uint `gorm:"type:int;not null;comment:书籍id"`
	CategoryID uint `gorm:"type:int;not null;comment:分类id"`
	Creator    uint `gorm:"type:int;not null;comment:创建者"`
}
