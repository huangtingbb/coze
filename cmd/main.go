package main

import (
	"coze-agent-platform/config"
	"coze-agent-platform/middleware"
	"coze-agent-platform/models"
	"coze-agent-platform/routers"
	"coze-agent-platform/utils"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "coze-agent-platform/docs" // 导入生成的文档
)

// @title Coze Agent Platform API
// @version 1.0
// @description 基于Gin框架的Coze Agent中台微服务系统
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志
	utils.InitLogger()

	// 初始化数据库
	models.InitDB()

	// 初始化Redis
	utils.InitRedis()

	// 设置Gin模式
	if config.Cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "service is running",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置路由
	routers.SetupRoutes(r)

	// 启动服务器
	port := ":" + config.Cfg.App.Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
