package api

import (
	"fmt"
	"gin_admin_system/internal/app"
	"gin_admin_system/internal/app/ginx"
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var LoginSet = wire.NewSet(wire.Struct(new(LoginAPI), "*"))

type LoginAPI struct {
	LoginSrv *service.LoginSrv
}

func (l *LoginAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var item types.LoginReq
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		// app.ReturnWithDetailed()
	}

	// 验证码的校验
	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		// 无效的验证码
		return
	}

	// 验证账号 密码
	userItem, err := l.LoginSrv.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		// 账号密码验证失败
		// app.ReturnWithDetailed()
	}

	// 生成token
	tokenInfo, err := l.LoginSrv.GenerateToken(ctx, l.formatTokenUserID(userItem.ID, userItem.UserName))
	if err != nil {
		// app.ReturnWithDetailed()
	}

	app.OkWithData(tokenInfo, c)
}

func (l *LoginAPI) formatTokenUserID(userID uint64, userName string) string {
	return fmt.Sprintf("%d-%s", userID, userName)
}

func (l *LoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	err := l.LoginSrv.DestroyToken(ctx, ginx.GetToken(c))
	if err != nil {
		//
		// app.ReturnWithDetailed()
	}
	app.Ok(c)
}
