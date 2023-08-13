package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DocRouteV1(c *gin.RouterGroup) {

	// 文档详情
	c.GET("/detail/:uuid", DocDetail)
	c.GET("/test", test)
	c.POST("/create/markdown", DocCreateMarkdown)
	c.POST("/create/upload", DocCreateUpload)
	c.GET("/get/:uuid", GetFile)

}

// DocCreateUpload 创建上传文档
// @Summary 创建上传文档
// @Description 创建上传文档 POST
// @Tags 文档
// @Accept
// @Produce json
// @Param book_uuid body string true "文档库uuid"
// @Param title body string true "文档标题"
// @Param file body string true "文档文件"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
func DocCreateUpload(c *gin.Context) {

	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	// 获取POST参数
	file, err := c.FormFile("files")
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	// 获取GET参数
	bookUUID := c.Query("book_uuid")

	// 获取文档库
	var book Book
	if err := DB.Where("uuid = ?", bookUUID).First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			BadRequest(c, "文档库不存在"+bookUUID)
			return
		}
		InternalServerError(c, err.Error())
		return
	}

	// 获取文件类型
	fileType := GetFileType(file.Filename)
	if fileType == "unknown" {
		BadRequest(c, "无法识别的文件类型")
		return
	}

	uuidStr := uuid.New().String()

	// 保存文件
	filePath := fmt.Sprintf("./files/%s/1/%s", uuidStr, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// 创建文档
	var doc Document
	doc.UUID = uuid.New().String()
	doc.Name = file.Filename
	doc.Type = fileType
	doc.Creator = userID
	doc.Path = filePath
	if err := DB.Create(&doc).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// 创建文档版本
	var docVersion DocumentVersion
	docVersion.DocumentID = doc.ID
	docVersion.Version = "1"
	docVersion.Creator = userID
	docVersion.Path = filePath
	docVersion.UUID = uuid.New().String()
	if err := DB.Create(&docVersion).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	// 创建对应
	var bookDoc DocumentList
	bookDoc.BookID = book.ID
	bookDoc.DocumentID = doc.ID
	bookDoc.CategoryID = 1

	if err := DB.Create(&bookDoc).Error; err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Success(c, gin.H{
		"uuid": doc.UUID,
	})

}

// GetFile
func GetFile(c *gin.Context) {

	// 鉴权
	JWTAuthMiddleware(c)

	Uuid := c.Param("uuid")

	// 查询这个文档
	var doc Document
	DB.Where("uuid = ?", Uuid).Find(&doc)

	if doc.ID == 0 {
		BadRequest(c, "不存在这个文档")
		return
	}

	// 获取版本
	var docVersion DocumentVersion
	DB.Where("document_id = ?", doc.ID).Order("created_at desc").First(&docVersion)

	if docVersion.ID == 0 {
		BadRequest(c, "无法找到文档版本")
		return
	}

	// 如果这个文档是markdown
	if doc.Type == "markdown" {
		// 打开 /files/{uuid}.md 文件并且返回
		filePath := fmt.Sprintf("./files/%s.%s.md", Uuid, docVersion.Version)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			InternalServerError(c, "无法读取文档内容"+filePath)
			return
		}

		Success(c, string(content))
		return
	}

	if (doc.Type == "docx") || (doc.Type == "xlsx") || (doc.Type == "pptx") {
		filepath := fmt.Sprintf("%s", docVersion.Path)
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			InternalServerError(c, "无法读取文档内容"+filepath)
			return
		}

		Success(c, string(content))
		return
	}
}

// DocCreateMarkdown 创建markdown文档
// @Summary 创建markdown文档
// @Description 创建markdown文档 POST
// @Tags 文档
// @Accept json
// @Produce json
// @Param name body string true "文档名称"
// @Param book_uuid body string true "文档库uuid"
// @Param content body string true "文档内容"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
func DocCreateMarkdown(c *gin.Context) {

	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	// 获取POST参数
	var req struct {
		Name      string `json:"name" binding:"required"`
		Cate      string `json:"category" binding:"required"`
		BookUUID  string `json:"book_uuid" binding:"required"`
		Content   string `json:"content" binding:"required"`
		IsPrivate bool   `json:"is_private"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 获取文档库
	var book Book
	if err := DB.Where("uuid = ?", req.BookUUID).First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			BadRequest(c, "文档库不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}

	var DocVersion DocumentVersion

	// 将markdown转换.md文件, 保存在/files/目录下
	// 生成uuid
	uuidStr := uuid.New().String()
	// 生成文件名, 版本号1
	fileName := uuidStr + ".1.md"
	// 生成文件路径
	filePath := "./files/" + fileName
	// 保存文件
	if err := saveFile(filePath, req.Content); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 创建
	doc := Document{
		Name:    req.Name,
		UUID:    uuidStr,
		Path:    filePath,
		Type:    "markdown",
		Creator: userID,
	}
	if err := DB.Create(&doc).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 创建版本
	DocVersion.Creator = userID
	DocVersion.UUID = uuid.New().String()
	DocVersion.DocumentID = doc.ID
	DocVersion.Version = "1"
	DocVersion.Path = filePath
	DocVersion.Size = int64(len(req.Content))

	shaStr, err := sha256File(filePath)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	DocVersion.SHA256 = shaStr

	if err := DB.Create(&DocVersion).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 搜索是否有这个分类， 没有则新建
	var category Category
	DB.Where("name = ?", req.Cate).Where("book_id = ?", book.ID).First(&category)

	if category.ID == 0 {
		category.Creator = userID
		category.Name = req.Cate
		category.UUID = uuid.New().String()
		category.BookID = book.ID

		DB.Save(&category)
	}

	// 创建文档库文档关联
	bookDoc := DocumentList{
		DocumentID: doc.ID,
		BookID:     book.ID,
		CategoryID: category.ID,
		Creator:    userID,
	}

	if err := DB.Create(&bookDoc).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 返回
	Success(c, doc)

}

func test(c *gin.Context) {

	// 获取GET参数
	var req struct {
		Type string `form:"type" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	switch req.Type {
	case "docx":
		// 返回/test/test.docx文件
		c.File("./test/test.docx")
	case "pdf":
		// 返回/test/test.pdf文件
		c.File("./test/test.pdf")
	case "md":
		// 返回/test/test.md文件
		c.File("./test/test.md")
	default:
		BadRequest(c, "不支持的类型")
	}

}

// DocDetail 文档详情
// @Summary 文档详情
// @Description 文档详情 GET
// @Tags 文档
// @Accept json
// @Produce json
// @Param uuid path string true "文档uuid"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
func DocDetail(c *gin.Context) {
	// 验证
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	uuidStr := c.Param("uuid")

	var doc Document
	if err := DB.Where("uuid = ?", uuidStr).First(&doc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			BadRequest(c, "文档不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}

	pub := uuid.New().String()
	token, err := GenerateJWTToken(pub)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 生成访问链接
	doc.URL = fmt.Sprintf("/api/v1/doc/get/%s?token=%s&pub=%s", doc.UUID, token, pub)

	Success(c, doc)

}
