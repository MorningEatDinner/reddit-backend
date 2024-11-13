package routes

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	gs "github.com/swaggo/gin-swagger"
	"github.com/xiaorui/reddit-async/reddit-backend/controller"
	"github.com/xiaorui/reddit-async/reddit-backend/logger"
	"github.com/xiaorui/reddit-async/reddit-backend/middlewares"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置为发布模式
	}
	r := gin.New() // 我尝试进行新的变化
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// r.Use(middlewares.RateLimitMiddleware(time.Second*2, 10)) // 2s新增1个令牌， 容量为10
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		// 创建验证相关的路由组
		authGroup := v1.Group("/auth")
		// 注册业务路由
		{
			// 这是旧的注册方式， 后面不使用了
			authGroup.POST("/signup", controller.SignUpHandler)
			authGroup.POST("/login", controller.LoginHandler)

			// 下面是新增的
			// 用户进行注册之前， 前端发送请求看是否已经注册， 如果已经注册则不能注册了
			authGroup.POST("/signup/phone/exist", controller.IsPhoneExist)
			authGroup.POST("/signup/email/exist", controller.IsEmailExist)

			// 验证码相关
			authGroup.GET("/code/captcha", controller.GetCaptcha)
			authGroup.POST("/code/phone", controller.SendPhoneCode)
			authGroup.POST("/code/email", controller.SendEmailCode)

			// 使用手机或者邮箱进行注册
			// 注册方式： 1. 手机+密码 2. 邮箱+密码
			authGroup.POST("/signup/phone", controller.SignupUsingPhone)
			authGroup.POST("/signup/email", controller.SignupUsingEmail)

			// 登录相关
			// 登录方式： 1. 手机+验证码， 2. 邮箱——密码
			authGroup.POST("/login/phone", controller.LoginUsingPhone)
			authGroup.POST("/login/email", controller.LoginUsingEmail)
			//authGroup.POST("/login/username", controller.LoginUsingUsername)
			authGroup.GET("/login/refresh-token", controller.RefreshToken)

			// 重置密码
			//authGroup.POST("/password/phone", nil)
			//authGroup.POST("/password/email", nil)
		}
		// 专门给周报使用的端点
		v1.GET("/week_report", controller.GetPostListHandler0)
		v1.GET("/all_email", controller.GetEmailList)
		v1.POST("/test_async", controller.TestAsync)
		v1.POST("/test_mq", controller.TestMq)

		// 后面的所有请求都需要使用这个中间件，即需要验证是否进行了登陆
		v1.Use(middlewares.JWTAuthMiddleware()) // 调用Use这个方法， 传入的中间件会被注入当下这个路由组中
		// 创建用户相关的路由组
		usersGroup := v1.Group("/user")
		{
			usersGroup.GET("", controller.CurrentUser)
			usersGroup.PUT("", controller.UpdateProfile) // 更新用户信息
			usersGroup.PUT("/email", controller.UpdateEmail)
			usersGroup.PUT("/phone", controller.UpdatePhone)
			usersGroup.PUT("/password", controller.UpdatePassword) // 更改密码
			usersGroup.PUT("/avatar", controller.UpdateAvatar)     // 更新头像
		}

		commGroup := v1.Group("/community")
		{
			commGroup.POST("", controller.CreateNewCommunity)        // 新建社区
			commGroup.GET("", controller.CommunityHandler)           //  获取所有社区信息
			commGroup.GET("/:id", controller.CommunityDetailHandler) // 获取当个社区的详细信息
			commGroup.PUT("/:id", controller.UpdateCommunity)        // 更新单个社区的信息
			commGroup.DELETE("/:id", controller.DeleteCommunity)     // 删除某个社区
		}

		postGroup := v1.Group("/post")
		{
			postGroup.POST("", controller.CreatePostHandler) // 创建帖子
			postGroup.GET("/:id", controller.GetPostHandler) // 获取某个具体帖子的信息
			//postGroup.GET("/posts", controller.GetPostListHandler)
			postGroup.GET("/posts2", controller.GetPostListHandler0) // 不定社区
			postGroup.POST("/vote", controller.PostVoteHandler)      // 对于某个帖子进行投票
			postGroup.GET("/posts3", controller.GetPostListHandler0) // 给定社区

			postGroup.DELETE("/:id", controller.DeletePost) // 删除删除

			commentGroup := postGroup.Group("/comment")
			{
				commentGroup.POST("/:post_id", controller.CreateComment)      // 给某个post发送一个comment
				commentGroup.GET("/:post_id", controller.GetComment)          // 获取某个post的所有comment
				commentGroup.DELETE("/:comment_id", controller.DeleteComment) // 删除某个comment
			}
		}
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	pprof.Register(r)

	return r
}
