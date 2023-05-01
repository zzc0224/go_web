package routers

import (
	"bluebell_backend/controller"
	"log"
	"net/http"
	"time"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	_ "bluebell_backend/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r := gin.Default()
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	v1.POST("/login", controller.LoginHandler)
	v1.POST("/signup", controller.SignUpHandler)
	v1.GET("/refresh_token", controller.RefreshTokenHandler)
	v1.GET("/post/:id", controller.PostDetailHandler)
	v1.GET("/post", controller.PostListHandler) //redis
	v1.Use(controller.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)

		v1.GET("/communityList", controller.CommunityListHandler) //根据社区查找帖子

		v1.GET("/post2", controller.PostList2Handler) //mysql

		v1.POST("/vote", controller.VoteHandler)

		v1.POST("/comment/:id", controller.CommentHandler)
		v1.GET("/comment/:id", controller.CommentListHandler)

		v1.GET("/recommend", controller.ReCommend)

		v1.GET("/CurrUserName", controller.GetCurrUserName)

		v1.POST("/upLoad", controller.UploadImg)

		v1.POST("/search", controller.Search)

		v1.GET("/Space", controller.GetPostListBYUser) //个人发的帖子

		v1.GET("/OtherUserSpace/:id", controller.GetOtherUserPost)

		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong test")
		})

	}
	pprof.Register(r)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
func Cors2() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			//Access-Control-Allow-Origin是必须的,他的值要么是请求Origin字段的值,要么是一个*, 表示接受任意域名的请求
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			//该字段可选。CORS请求时，XMLHttpRequest对象的getResponseHeader()方法只能拿到6个基本字段：Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma。
			//如果想拿到其他字段，就必须在Access-Control-Expose-Headers里面指定。上面的例子指定，getResponseHeader('FooBar')可以返回FooBar字段的值。
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			//该字段可选，用来指定本次预检请求的有效期，单位为秒。有效期是20天（1728000秒），即允许缓存该条回应1728000秒（即20天），在此期间，不用发出另一条预检请求。
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func Cors3() gin.HandlerFunc {
	c := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Access-Token", "Authorization"},
		MaxAge:          6 * time.Hour,
	}

	return cors.New(c)
}
