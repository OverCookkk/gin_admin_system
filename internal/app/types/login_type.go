package types

type LoginReq struct {
	UserName    string `json:"user_name" binding:"required"`    // 用户名
	Password    string `json:"password" binding:"required"`     // 密码(md5加密)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}

type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	// TokenType   string `json:"token_type"`   // 令牌类型
	// ExpiresAt   int64  `json:"expires_at"`   // 过期时间戳
}
