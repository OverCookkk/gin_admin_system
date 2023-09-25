package app

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/middleware"
	"gin_admin_system/internal/app/router"
	"github.com/gin-gonic/gin"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	// TODO:中间件
	// todo:AllowPathPrefixSkipper

	// Recover：painc时候记录堆栈信息
	app.Use(middleware.RecoveryMiddleware())

	// trace ID
	app.Use(middleware.TraceMiddleware())

	// CORS跨域配置
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	r.Register(app)

	return app
}
