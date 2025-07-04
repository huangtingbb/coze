package routers

import (
	"coze-agent-platform/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// API路由组
	api := r.Group("/api")

	// 公开路由（无需认证）
	public := api.Group("/")
	{
		public.POST("/auth/login", controllers.Login)
		public.POST("/auth/register", controllers.Register)

		api.Group("/coze").GET("/token", controllers.GetCozeToken)
	}

	// 需要认证的路由
	auth := api.Group("/")
	// auth.Use(middleware.JWTAuth())
	{
		// 用户相关
		auth.GET("/users/profile", controllers.GetUserProfile)
		auth.PUT("/users/profile", controllers.UpdateUserProfile)

		// Agent相关
		auth.GET("/agents", controllers.ListAgents)
		auth.POST("/agents", controllers.CreateAgent)
		auth.GET("/agents/:id", controllers.GetAgent)
		auth.PUT("/agents/:id", controllers.UpdateAgent)
		auth.DELETE("/agents/:id", controllers.DeleteAgent)

		// 对话相关
		auth.GET("/conversations", controllers.ListConversations)
		auth.POST("/conversations", controllers.CreateConversation)
		auth.GET("/conversations/:id", controllers.GetConversation)
		auth.DELETE("/conversations/:id", controllers.DeleteConversation)

		// 消息相关
		auth.GET("/conversations/:id/messages", controllers.GetMessages)
		auth.POST("/conversations/:id/messages", controllers.SendMessage)
		auth.POST("/conversations/messages/stream", controllers.SendMessageStream)

		// 文件上传
		auth.POST("/common/upload/file", controllers.UploadFile)
	}
}
