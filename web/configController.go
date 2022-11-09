package web

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/johlanse/study_xxqg/conf"
)

func configFileGet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		level := ctx.GetInt("level")
		if level != 1 {
			ctx.JSON(200, Resp{
				Code:    403,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "获取成功",
			Data:    conf.GetConfigFile(),
			Success: true,
			Error:   "",
		})

	}
}

func configFileSet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		level := ctx.GetInt("level")
		if level != 1 {
			ctx.JSON(200, Resp{
				Code:    403,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			return
		}
		var body map[string]string
		_ = ctx.ShouldBindJSON(&body)

		err := yaml.Unmarshal([]byte(body["data"]), new(conf.Config))
		if err != nil {
			ctx.JSON(200, Resp{
				Code:    503,
				Message: "配置提交失败！！",
				Data:    nil,
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		err = conf.SaveConfigFile(body["data"])
		if err != nil {
			ctx.JSON(200, Resp{
				Code:    503,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "获取成功",
			Data:    conf.GetConfigFile(),
			Success: true,
			Error:   "",
		})

	}
}

func configGet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		level := ctx.GetInt("level")
		if level != 1 {
			ctx.JSON(200, Resp{
				Code:    403,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "获取成功",
			Data:    conf.GetConfig(),
			Success: true,
			Error:   "",
		})

	}
}

func configSet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		level := ctx.GetInt("level")
		if level != 1 {
			ctx.JSON(200, Resp{
				Code:    403,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   "",
			})
			return
		}
		c := new(conf.Config)
		err := ctx.BindJSON(c)
		if err != nil {
			ctx.JSON(200, Resp{
				Code:    401,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		err = conf.SetConfig(*c)
		if err != nil {
			ctx.JSON(200, Resp{
				Code:    503,
				Message: "",
				Data:    nil,
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(200, Resp{
			Code:    200,
			Message: "",
			Data:    nil,
			Success: true,
			Error:   "",
		})
	}
}
