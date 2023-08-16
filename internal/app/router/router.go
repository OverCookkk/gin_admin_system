package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
	// Auth jwt验证
	// Casbin 权限控制
	// LoginAPI *api

}

func (r *Router) Register(app *gin.Engine) error {
	return nil
}

func (r *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

// func (r *Router) RegisterAPI(app *gin.Engine) {
// 	g := app.Group("/api")
// }
