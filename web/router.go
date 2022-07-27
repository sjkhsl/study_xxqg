package web

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/huoxue1/study_xxqg/conf"
	"github.com/huoxue1/study_xxqg/utils"
)

//go:embed xxqg/build
var static embed.FS

// RouterInit
// @Description:
// @return *gin.Engine
func RouterInit() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())

	router.StaticFS("/static", http.FS(static))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/static/xxqg/build/home.html")
	})

	auth := router.Group("/auth")
	auth.POST("/login", Login())
	auth.POST("/check/:token", CheckToken())

	if utils.FileIsExist("dist") {
		router.StaticFS("/dist", gin.Dir("./dist/", true))
	}

	user := router.Group("/user", check())
	// 添加用户
	user.POST("", addUser())

	user.GET("/", getUsers())

	router.GET("/score", getScore())

	router.POST("/study", study())

	router.POST("/stop_study", check(), stopStudy())

	router.GET("/log", check(), getLog())

	router.GET("/sign/*proxyPath", check(), sign())
	router.GET("/login/*proxyPath", check(), generate())
	router.POST("/login/*proxyPath", check(), generate())
	return router
}

func check() gin.HandlerFunc {
	config := conf.GetConfig()
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("xxqg_token")
		if token == "" || (utils.StrMd5(config.Web.Account+config.Web.Password) != token) {
			ctx.JSON(403, Resp{
				Code:    403,
				Message: "the auth fail",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
