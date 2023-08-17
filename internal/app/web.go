package app

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/router"
	"github.com/gin-gonic/gin"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	// 中间件

	r.Register(app)

	return app
}
