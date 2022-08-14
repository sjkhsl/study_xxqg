// Package web
// @Description: 封装了所以web相关的内容
package web

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/utils"
)

// 将静态文件嵌入到可执行程序中来
//go:embed xxqg/build
var static embed.FS

//go:embed app
var newUI embed.FS

// RouterInit
// @Description:
// @return *gin.Engine
func RouterInit() *gin.Engine {
	router := gin.Default()
	router.RemoveExtraSlash = true
	router.Use(cors())

	// 挂载静态文件
	router.StaticFS("/static", http.FS(static))
	// 访问首页时跳转到对应页面
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/static/xxqg/build/home.html")
	})

	router.GET("/about", func(context *gin.Context) {
		context.JSON(200, utils.GetAbout())
	})

	router.GET("/new/*file", func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.Path, "js") {
			ctx.Header("Content-Type", "application/javascript; charset=utf-8")
		}
		f, err := newUI.ReadFile(strings.TrimLeft(ctx.Param("file"), "/"))
		if err != nil {
			log.Errorln(err.Error())
			f, _ = newUI.ReadFile("app/index.html")
		}
		_, _ = ctx.Writer.Write(f)
		ctx.Status(200)
	})
	// 对权限的管理组
	auth := router.Group("/auth")
	// 用户登录的接口
	auth.POST("/login", userLogin())
	// 检查登录状态的token是否正确
	auth.POST("/check/:token", checkToken())

	// 对于用户可自定义挂载文件的目录
	if utils.FileIsExist("./config/dist/") {
		router.StaticFS("/dist", http.Dir("./config/dist/"))
	}

	// 对用户管理的组
	user := router.Group("/user", check())
	// 添加用户
	user.POST("", addUser())
	// 获取所以已登陆的用户
	user.GET("", getUsers())
	// 删除用户
	user.DELETE("", deleteUser())

	// 获取用户成绩
	router.GET("/score", getScore())
	// 让一个用户开始学习
	router.POST("/study", study())
	// 让一个用户停止学习
	router.POST("/stop_study", check(), stopStudy())
	// 获取程序当天的运行日志
	router.GET("/log", check(), getLog())

	// 登录xxqg的三个接口
	router.GET("/sign/", sign())
	router.GET("/login/*proxyPath", generate())
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
