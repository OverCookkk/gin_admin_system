package api

import (
	"fmt"
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/contextx"
	"gin_admin_system/internal/app/ginx"
	"gin_admin_system/internal/app/response"
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var LoginApiSet = wire.NewSet(wire.Struct(new(LoginAPI), "*"))

type LoginAPI struct {
	LoginSrv *service.LoginSrv
}

func (l *LoginAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var item types.LoginReq
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}

	// 验证码的校验
	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) { // CaptchaID为后端返回的id，CaptchaCode为用户输入的验证码
		// 无效的验证码
		response.JsonError(c, errors.New("Invalid verification code"))
		return
	}

	// 验证账号 密码
	userItem, err := l.LoginSrv.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		// 账号密码验证失败
		response.JsonError(c, err)
		return
	}

	// 生成token
	tokenInfo, err := l.LoginSrv.GenerateToken(ctx, l.formatTokenUserID(userItem.ID, userItem.UserName))
	if err != nil {
		response.JsonError(c, err)
		return
	}

	response.JsonData(c, tokenInfo)
}

func (l *LoginAPI) formatTokenUserID(userID uint64, userName string) string {
	return fmt.Sprintf("%d-%s", userID, userName)
}

func (l *LoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	err := l.LoginSrv.DestroyToken(ctx, ginx.GetToken(c))
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

// GetCaptcha 获取验证码id
func (l *LoginAPI) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()

	item, err := l.LoginSrv.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, item)
}

// ResCaptcha 通过GetCaptcha接口返回的id获取验证码图片
func (l *LoginAPI) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		response.JsonError(c, errors.New("captcha id not empty"))
		return
	}

	if c.Query("reload") != "" { // 需要重新获取新的captcha_id
		if !captcha.Reload(captchaID) {
			response.JsonError(c, errors.New("not found captcha id"))
			return
		}
	}

	// 验证码图片写入响应c.Writer中
	err := l.LoginSrv.ResCaptcha(ctx, c.Writer, captchaID, config.C.Captcha.Width, config.C.Captcha.Height)
	if err != nil {
		response.JsonError(c, err)
	}
}

func (l *LoginAPI) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	userInfo, err := l.LoginSrv.GetLoginInfo(ctx, contextx.GetUserID(ctx))
	if err != nil {
		response.JsonError(c, err)
		return
	}

	response.JsonData(c, userInfo)
}
