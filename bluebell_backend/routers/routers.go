package routers

import (
	"bluebell_backend/controller"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	_ "bluebell_backend/docs" // 千万不要忘了导入把你上一步生成的docs

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
