package main

import "github.com/gin-gonic/gin"

func BadRequest(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 400,
		"data": nil,
		"msg":  msg,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
		"msg":  "ok",
	})
}

func SuccessMsg(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": nil,
		"msg":  msg,
	})
}

func InternalServerError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 500,
		"data": nil,
		"msg":  msg,
	})
}
