package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// google uuid

func BookRouteV1(c *gin.RouterGroup) {
	c.GET("/info/:uuid", BookDetail)
	c.POST("/create", BookCreate)
	c.POST("/category/search", BookCatalogSearch)
	c.GET("/list", GetBookList)
}

// GetBookList
func GetBookList(c *gin.Context) {
	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	var books []Book
	DB.Find(&books)

	Success(c, books)

}

// BookCatalogSearch 搜索文档库目录
// @Summary 搜索文档库目录
// @Description 搜索文档库目录 POST
// @Tags 文档库
// @Accept json
// @Produce json
// @Param book_uuid body string true "文档库uuid"
// @Param keyword body string true "关键词"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
func BookCatalogSearch(c *gin.Context) {
	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	// 获取POST参数
	var req struct {
		BookUUID string `json:"book_uuid" binding:"required"`
		Keyword  string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 获取文档库
	var book Book
	if err := DB.Where("uuid = ?", req.BookUUID).First(&book).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 获取文档库目录
	var catalog []Category
	if err := DB.Where("book_id = ? AND name LIKE ?", book.ID, "%"+req.Keyword+"%").Find(&catalog).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回
	Success(c, catalog)
}

// BookCreate 创建文档库
// @Summary 创建文档库
// @Description 创建文档库 POST
// @Tags 文档库
// @Accept json
// @Produce json
// @Param name body string true "文档库名称"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
func BookCreate(c *gin.Context) {
	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	// 获取POST参数
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 创建
	book := Book{
		Name: req.Name,
		UUID: uuid.New().String(),
	}
	if err := DB.Create(&book).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回
	Success(c, book)
}

// BookDetail 文档库详情
// @Summary 文档库详情
// @Description 文档库详情 GET
// @Tags 文档库
// @Accept json
// @Produce json
// @Param uuid path string true "文档库uuid"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
func BookDetail(c *gin.Context) {
	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	// 获取uuid
	bookuuid := c.Param("uuid")

	// 查询
	var book Book
	if err := DB.Where("uuid = ?", bookuuid).First(&book).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 获取所有文档
	var docslist []DocumentList
	if err := DB.Where("book_id = ?", book.ID).Find(&docslist).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 查询出全部分类的名字
	var category []Category
	if err := DB.Where("book_id = ?", book.ID).Find(&category).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	type docsData struct {
		Cate      string
		UUID      string
		Documents []struct {
			Title string
			UUID  string
		}
	}

	var data []docsData

	for _, v := range category {
		var docs docsData
		docs.Cate = v.Name
		docs.UUID = v.UUID
		for _, v2 := range docslist {
			if v2.CategoryID == v.ID {
				var doc Document
				if err := DB.Where("id = ?", v2.DocumentID).First(&doc).Error; err != nil {
					BadRequest(c, err.Error())
					return
				}
				docs.Documents = append(docs.Documents, struct {
					Title string
					UUID  string
				}{Title: doc.Name, UUID: doc.UUID})
			}
		}
		data = append(data, docs)
	}

	var res struct {
		BookInfo Book
		Category []Category
		Docs     []docsData
	}

	res.BookInfo = book
	res.Category = category
	res.Docs = data

	// 返回
	Success(c, res)

}
