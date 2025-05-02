package handler

import (
	"awesomeProject/internal/handler/app"
	"awesomeProject/internal/handler/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {

	app.InitApp()

	// 注册中间件
	r.Use(middleware.CORS())

	r.GET("/hello", helloHandler)
	r.GET("/refresh", middleware.RefreshTokenHandler)

	r.GET("/check", app.CheckAnswer)

	// 用户认证相关路由
	r.POST("/api/v1/register", app.RegisterHandler) // 用户注册
	r.POST("/api/v1/login", app.LoginHandler)       // 用户登录
	r.POST("/api/v1/logout", app.LogoutHandler)     // 用户注销

	// 需要认证的路由组
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/user", app.GetCurrentUserHandler) // 获取当前用户信息
	}

	// 课程相关路由 (根据API文档，这些路由不需要认证)
	courseGroup := r.Group("/api/v1/courses")
	{
		courseGroup.GET("", app.GetAllCoursesHandler)                             // 获取所有课程
		courseGroup.GET("/:course_id", app.GetCourseByIDHandler)                  // 获取特定课程信息
		courseGroup.GET("/:course_id/references", app.GetCourseReferencesHandler) // 获取课程参考书
		courseGroup.GET("/experiment-status", app.GetCourseExperimentStatusHandler)
	}

}

func helloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
