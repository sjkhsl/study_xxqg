package web

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed xxqg/build
var static embed.FS

// RouterInit
// @Description:
// @return *gin.Engine
func RouterInit() *gin.Engine {
	router := gin.New()
	router.Use(Cors())

	router.StaticFS("/static", http.FS(static))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/static/xxqg/build/home.html")
	})
	user := router.Group("/user")
	// 添加用户
	user.POST("/", addUser())

	user.GET("/", getUsers())

	router.GET("/score", getScore())

	router.POST("/study", study())

	router.POST("/stop_study", stopStudy())

	router.GET("/log", getLog())

	router.GET("/sign/*proxyPath", sign())
	router.GET("/login/*proxyPath", generate())
	router.POST("/login/*proxyPath", generate())
	return router
}
