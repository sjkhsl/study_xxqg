// Package web
// @Description:
package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/lib/state"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/push"
	"github.com/johlanse/study_xxqg/utils"
)

// checkToken
/* @Description:
 * @return gin.HandlerFunc
 */
func checkToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Param("token")
		config := conf.GetConfig()
		md5 := utils.StrMd5(config.Web.Account + config.Web.Password)
		if md5 == token {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "",
				Data:    1,
				Success: true,
				Error:   "",
			})
		} else if checkCommonUser(token) {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "",
				Data:    2,
				Success: true,
				Error:   "",
			})
		} else {
			ctx.JSON(200, Resp{
				Code:    403,
				Message: "",
				Data:    -1,
				Success: false,
				Error:   "",
			})
		}
	}
}

func userLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type user struct {
			Account  string `json:"account"`
			Password string `json:"password"`
		}
		u := new(user)
		_ = ctx.BindJSON(u)
		config := conf.GetConfig()
		if u.Account == config.Web.Account && u.Password == config.Web.Password {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "登录成功，尊贵的管理员用户",
				Data:    utils.StrMd5(u.Account + u.Password),
				Success: true,
				Error:   "",
			})
		} else if checkCommonUser(utils.StrMd5(u.Account + u.Password)) {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "登录成功",
				Data:    utils.StrMd5(u.Account + u.Password),
				Success: true,
				Error:   "",
			})
		} else {
			ctx.JSON(200, Resp{
				Code:    403,
				Message: "登录失败，请联系管理员",
				Data:    "",
				Success: false,
				Error:   "",
			})
		}
	}
}

func getScore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		score, err := lib.GetUserScore(model.TokenToCookies(token))
		if err != nil {
			ctx.JSON(403, Resp{
				Code:    403,
				Message: "",
				Data:    err.Error(),
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "获取成功",
			Data:    lib.FormatScore(score),
			Success: true,
			Error:   "",
		})
	}
}

func addUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type params struct {
			Code  string `json:"code"`
			State string `json:"state"`
		}
		p := new(params)
		err := ctx.BindJSON(p)
		if err != nil {
			ctx.JSON(403, Resp{
				Code:    403,
				Message: "",
				Data:    err.Error(),
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		_, err = lib.GetToken(p.Code, p.State, ctx.GetString("token"))
		if err != nil {
			ctx.JSON(403, Resp{
				Code:    403,
				Message: "",
				Data:    err.Error(),
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "登录成功",
			Data:    "登录成功",
			Success: true,
			Error:   "",
		})
	}
}

func getExpiredUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		failUser, err := model.QueryFailUser()
		if err != nil {
			nilArray := make([]interface{}, 0)
			if err == sql.ErrNoRows {
				ctx.JSON(200, Resp{
					Code:    200,
					Message: "",
					Data:    nilArray,
					Success: true,
					Error:   "",
				})
			} else {
				ctx.JSON(502, Resp{
					Code:    502,
					Message: "",
					Data:    nilArray,
					Success: false,
					Error:   err.Error(),
				})
			}
			return
		}
		level := ctx.GetInt("level")
		if level == 1 {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "",
				Data:    failUser,
				Success: true,
				Error:   "",
			})
		} else {
			var myFaileUser []*model.User
			for _, user := range failUser {
				if user.Token == ctx.GetString("token") {
					myFaileUser = append(myFaileUser, user)
				}
			}
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "",
				Data:    myFaileUser,
				Success: true,
				Error:   "",
			})
		}

	}
}

func getUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := model.Query()
		if err != nil {
			return
		}
		if users == nil {
			ctx.JSON(200, Resp{
				Code:    200,
				Message: "查询成功",
				Data:    []interface{}{},
				Success: true,
				Error:   "",
			})
			return
		}
		level := ctx.GetInt("level")
		if level != 1 {
			users, err = model.QueryByPushID(ctx.GetString("token"))
			if err != nil {
				return
			}
			if users == nil {
				ctx.JSON(200, Resp{
					Code:    200,
					Message: "查询成功",
					Data:    []interface{}{},
					Success: true,
					Error:   "",
				})
				return
			}
		}

		var datas []map[string]interface{}
		for _, user := range users {
			datas = append(datas, map[string]interface{}{
				"nick":       user.Nick,
				"uid":        user.Uid,
				"token":      user.Token,
				"login_time": user.LoginTime,
				"is_study":   state.IsStudy(user.Uid),
			})
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "查询成功",
			Data:    datas,
			Success: true,
			Error:   "",
		})
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin != "" {
			// 接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// 允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "*")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "*")
			// 设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			// 允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 允许类型校验
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

func study() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("uid")
		user := model.Find(uid)
		core := &lib.Core{
			ShowBrowser: conf.GetConfig().ShowBrowser,
			Push:        push.GetPush(conf.GetConfig()),
		}
		core.Init()
		state.Add(user.Uid, core)
		config := conf.GetConfig()
		go func() {
			core.LearnArticle(user)
			core.LearnVideo(user)
			if config.Model == 2 {
				core.RespondDaily(user, "daily")
			} else if config.Model == 3 {
				core.RespondDaily(user, "daily")
				core.RespondDaily(user, "weekly")
				core.RespondDaily(user, "special")
			}
			state.Delete(uid)
		}()
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "",
			Data:    "",
			Success: true,
			Error:   "",
		})
	}
}

func stopStudy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("uid")
		core := state.Get(uid)
		core.Quit()
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "",
			Data:    "",
			Success: true,
			Error:   "",
		})
	}
}

func getLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.File(fmt.Sprintf("./config/logs/%v.log", time.Now().Format("2006-01-02")))
	}
}

func sign() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := utils.GetClient().R().Get("https://pc-api.xuexi.cn/open/api/sns/sign")
		if err != nil {
			ctx.JSON(403, Resp{
				Code:    403,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		data := response.Bytes()
		log.Debugln("访问sign结果返回内容为 ==》 " + string(data))
		ctx.Writer.WriteHeader(200)
		ctx.Writer.Write(data)
	}
}

func generate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		remote, _ := url.Parse("https://login.xuexi.cn/")
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = ctx.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = ctx.Param("proxyPath")
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// 删除用户
func deleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Query("uid")
		level := ctx.GetInt("level")
		if level != 1 {
			ctx.JSON(200, Resp{
				Code:    401,
				Message: "你没有权限删除用户！",
				Data:    "",
				Success: false,
				Error:   "你没有权限删除用户！",
			})
			return
		}
		err := model.DeleteUser(uid)
		if err != nil {
			ctx.JSON(200, Resp{
				Code:    503,
				Message: "",
				Data:    "",
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "删除成功",
			Data:    "",
			Success: true,
			Error:   "",
		})
	}
}
