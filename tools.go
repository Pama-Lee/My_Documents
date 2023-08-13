package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

// 计算字符串md5值
func Md5(String string) string {
	hash := md5.Sum([]byte(String))
	return hex.EncodeToString(hash[:])
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

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(c *gin.Context) {
	// 获取 Authorization header
	tokenStringFromHeader := c.GetHeader("Authorization")
	pubStringFromHeader := c.GetHeader("Pub")
	tokenStringFromQuery := c.Query("token")
	pubStringFromQuery := c.Query("pub")

	// 验证token格式
	if tokenStringFromHeader == "" && tokenStringFromQuery == "" {
		BadRequest(c, "请求未携带token，无权限访问")
		c.Abort()
		return
	}

	// 验证pub格式
	if pubStringFromHeader == "" && pubStringFromQuery == "" {
		BadRequest(c, "请求未携带pub，无权限访问")
		c.Abort()
		return
	}

	// 优先从header中获取token
	tokenString := tokenStringFromHeader
	if tokenStringFromHeader == "" {
		tokenString = tokenStringFromQuery
	}

	// 优先从header中获取pub
	pubString := pubStringFromHeader
	if pubStringFromHeader == "" {
		pubString = pubStringFromQuery
	}

	// 获取token
	tokenString = tokenString[7:]

	// 解析token
	token, claims, err := ParseJWTToken(tokenString)
	if err != nil || !token.Valid {
		BadRequest(c, "无效的token")
		c.Abort()
		return
	}

	// 验证通过后获取claim中的pub
	pub := (*claims)["pub"].(string)

	// 验证pub
	if pub != pubString {
		BadRequest(c, "无效的鉴权")
		c.Abort()
		return
	}

	c.Next()
}

// ParseJWTToken 解析JWTToken
func ParseJWTToken(tokenString string) (*jwt.Token, *jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.Jwt_secret), nil
	})

	// 解析错误
	if err != nil {
		return nil, nil, err
	}

	// 解析成功
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, err
}

// GenerateJWTToken 生成JWTToken
func GenerateJWTToken(pub string) (string, error) {
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pub": pub,
	})

	// 生成token字符串
	tokenString, err := token.SignedString([]byte(Config.Jwt_secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
