package router

import (
	"gin_admin_system/internal/app/api"
	"gin_admin_system/internal/app/middleware"
	"gin_admin_system/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// 编译前验证Router是否实现了IRouter的全部接口
var _ IRouter = (*Router)(nil)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
	Auth auth.JWTAuth // jwt验证
	// Casbin 权限控制
	// LoginAPI *api
	MenuApi *api.MenuApi
}

func (r *Router) Register(app *gin.Engine) error {
	return nil
}

func (r *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

func (r *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	// 中间件
	// todo:AllowPathPrefixSkipper
	g.Use(middleware.UserAuthMiddleware(r.Auth))

	v1 := g.Group("/v1")
	{
		// 菜单
		gMenu := v1.Group("/menus")
		{
			gMenu.GET("/menus", r.MenuApi.Query)
			gMenu.POST("Create", r.MenuApi.Create)
			gMenu.GET(":id", r.MenuApi.Get)
		}
	}
}
