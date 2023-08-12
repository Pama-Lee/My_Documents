package main

import "github.com/gin-gonic/gin"

func InitRouteV1(r *gin.RouterGroup) {

	// 用户路由
	user := r.Group("/user")
	UserRouteV1(user)

	// 文档路由
	doc := r.Group("/doc")
	DocRouteV1(doc)

	// 文档库路由
	book := r.Group("/book")
	BookRouteV1(book)

}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 允许所有跨域请求
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
