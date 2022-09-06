package web

import (
	"github.com/gin-gonic/gin"

	"github.com/johlanse/study_xxqg/conf"
)

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
