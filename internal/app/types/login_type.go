package types

type LoginReq struct {
	UserName    string `json:"user_name" binding:"required"`    // 用户名
	Password    string `json:"password" binding:"required"`     // 密码(md5加密)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}
