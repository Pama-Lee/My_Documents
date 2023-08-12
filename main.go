package main

// gin
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// gin-swagger UI
import "github.com/swaggo/gin-swagger/swaggerFiles"

// swagger 文档基本信息
// @title gin-swagger demo
// @description gin-swagger demo
// @version 1
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 打印banner
	printBanner()

	InitDB()

	// 获取.env配置
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// 判断是否为生产环境
	if gin.Mode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if os.Getenv("env") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化gin
	r := gin.Default()

	// 允许localhost跨域
	r.Use(Cors())

	// gin-swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	InitRouteV1(r.Group("/api/v1"))

	// 启动服务
	r.Run(":8080")
}

func printBanner() {
	fmt.Println(`

 ██████╗██╗   ██╗████████╗ ██████╗ ███╗   ██╗ ██████╗ 
██╔════╝╚██╗ ██╔╝╚══██╔══╝██╔═══██╗████╗  ██║██╔═══██╗
██║      ╚████╔╝    ██║   ██║   ██║██╔██╗ ██║██║   ██║
██║       ╚██╔╝     ██║   ██║   ██║██║╚██╗██║██║   ██║
╚██████╗   ██║      ██║   ╚██████╔╝██║ ╚████║╚██████╔╝
 ╚═════╝   ╚═╝      ╚═╝    ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝
		`)

	/**

	  ___  ___       ______                                      _
	  |  \/  |       |  _  \                                    | |
	  | .  . |_   _  | | | |___   ___ _   _ _ __ ___   ___ _ __ | |_ ___
	  | |\/| | | | | | | | / _ \ / __| | | | '_ ` _ \ / _ \ '_ \| __/ __|
	  | |  | | |_| | | |/ / (_) | (__| |_| | | | | | |  __/ | | | |_\__ \
	  \_|  |_/\__, | |___/ \___/ \___|\__,_|_| |_| |_|\___|_| |_|\__|___/
	           __/ |
	          |___/

	*/
	fmt.Println("\n___  ___       ______                                      _       \n|  \\/  |       |  _  \\                                    | |      \n| .  . |_   _  | | | |___   ___ _   _ _ __ ___   ___ _ __ | |_ ___ \n| |\\/| | | | | | | | / _ \\ / __| | | | '_ ` _ \\ / _ \\ '_ \\| __/ __|\n| |  | | |_| | | |/ / (_) | (__| |_| | | | | | |  __/ | | | |_\\__ \\\n\\_|  |_/\\__, | |___/ \\___/ \\___|\\__,_|_| |_| |_|\\___|_| |_|\\__|___/\n         __/ |                                                     \n        |___/                                                      \n")
}
