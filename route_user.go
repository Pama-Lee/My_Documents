package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func UserRouteV1(g *gin.RouterGroup) {
	// 用户登陆
	g.POST("/login", UserLogin)
	g.GET("/info", UserInfo)

}

// UserInfo 用户信息
// @Summary 用户信息
// @Description 用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
// @Router /user/info [get]
func UserInfo(c *gin.Context) {
	// 检验Header中的辨识信息
	userID := checkLoginInfo(c)
	if userID == 0 {
		BadRequest(c, "请先登陆")
		return
	}

	// 查询用户信息
	var user User
	if err := DB.Where("id = ?", userID).First(&user).Error; err != nil {
		BadRequest(c, err.Error())
		return
	}
	user.Password = ""

	// 返回
	Success(c, user)
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLogin 用户登陆
// @Summary 用户登陆
// @Description 用户登陆
// @Tags 用户
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 如果是测试环境
	if os.Getenv("env") == "debug" {
		if req.Username != "admin" || req.Password != "123456" {
			BadRequest(c, "用户名或密码错误")
			return
		} else {
			Success(c, gin.H{
				"token": "test_token",
			})
			return
		}
	} else {
		// 查询用户
		var user User
		if err := DB.Where("username = ? & password = ?", req.Username, req.Password).First(&user).Error; err != nil {
			BadRequest(c, "用户名或密码错误")
			return
		}

		// 生成token
		token := Md5(time.Now().String()+"_"+req.Username) + Md5(time.Now().String()+"_"+req.Password) + Md5(time.Now().String()+"_"+req.Username+"_"+req.Password)
		tokenEntity := Token{
			UserID: user.ID,
			Token:  token,
			Expire: time.Now().Unix() + 3600*24*7,
			Status: 1,
		}
		DB.Create(&tokenEntity)

		// 返回
		Success(c, token)
		return
	}
}

// 验证登陆信息
func checkLoginInfo(c *gin.Context) uint {
	token := c.GetHeader("Identify")
	if token == "" {
		return 0
	} else {

		// 如果是测试环境
		if os.Getenv("env") == "debug" {
			if token != "test_token" {
				return 0
			} else {
				return 1
			}
		}

		// 查询
		var tokenEntity Token
		if err := DB.Where("token = ?", token).Where("status = 1").First(&tokenEntity).Error; err != nil {
			return 0
		}
		// 检查是否过期
		nowTimestamp := time.Now().Unix()
		if nowTimestamp > tokenEntity.Expire {
			tokenEntity.Status = 0
			DB.Save(&tokenEntity)
			return 0
		} else {
			return tokenEntity.UserID
		}
	}
}
