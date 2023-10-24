package router

import (
	"gin_admin_system/internal/app/api"
	"gin_admin_system/internal/app/middleware"
	"gin_admin_system/pkg/auth"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// 编译前验证Router是否实现了IRouter的全部接口
var _ IRouter = (*Router)(nil)

// 也可以改成这样写，NewRouter返回Router实例，wire.Bind绑定实现与接口
// var RouterSet = wire.NewSet(NewRouter, wire.Bind(new(IRouter), new(*Router)))
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
	Auth           auth.JWTAuth           // Auth jwt验证
	CasbinEnforcer *casbin.SyncedEnforcer // Casbin 权限控制
	LoginApi       *api.LoginAPI
	MenuApi        *api.MenuApi
	RoleApi        *api.RoleApi
	UserApi        *api.UserApi
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

	g.Use(middleware.CasbinMiddleware(r.CasbinEnforcer))

	// todo:限流中间件

	v1 := g.Group("/v1")
	{
		// 登录相关
		gLogin := v1.Group("login")
		{
			gLogin.GET("captchaid", r.LoginApi.GetCaptcha)
			gLogin.GET("captcha", r.LoginApi.ResCaptcha)
			gLogin.POST("", r.LoginApi.Login)
			gLogin.POST("exit", r.LoginApi.Logout)
		}

		gCurrent := v1.Group("current")
		{
			// gCurrent.PUT("password", r.LoginApi.UpdatePassword)
			gCurrent.GET("user", r.LoginApi.GetUserInfo)
			// gCurrent.GET("menutree", r.LoginApi.QueryUserMenuTree)
		}
		// v1.POST("/refresh-token", r.LoginApi.RefreshToken)

		// 菜单
		gMenu := v1.Group("/menus")
		{
			gMenu.GET("", r.MenuApi.Query)
			gMenu.POST("Create", r.MenuApi.Create)
			gMenu.GET(":id", r.MenuApi.Get)
			gMenu.PUT(":id", r.MenuApi.Update)
			gMenu.DELETE(":id", r.MenuApi.Delete)
			gMenu.PATCH(":id/enable", r.MenuApi.Enable)
			gMenu.PATCH(":id/disable", r.MenuApi.Disable)
		}

		// 角色
		gRole := v1.Group("roles")
		{
			gRole.GET("", r.RoleApi.Query)
			gRole.GET(":id", r.RoleApi.Get)
			gRole.POST("", r.RoleApi.Create)
			gRole.PUT(":id", r.RoleApi.Update)
			gRole.DELETE(":id", r.RoleApi.Delete)
			gRole.PATCH(":id/enable", r.RoleApi.Enable)
			gRole.PATCH(":id/disable", r.RoleApi.Disable)
		}
		// v1.GET("/roles.select", r.RoleApi.QuerySelect)

		// 用户
		gUser := v1.Group("users")
		{
			gUser.GET("", r.UserApi.Query)
			gUser.GET(":id", r.UserApi.Get)
			gUser.POST("", r.UserApi.Create)
			gUser.PUT(":id", r.UserApi.Update)
			gUser.DELETE(":id", r.UserApi.Delete)
			gUser.PATCH(":id/enable", r.UserApi.Enable)
			gUser.PATCH(":id/disable", r.UserApi.Disable)
		}
	}
}
