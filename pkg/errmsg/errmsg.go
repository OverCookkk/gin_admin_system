package errmsg

const (
	SUCCSE = 200
	ERROR  = 500

	// code = 1000... 用户模块的错误
	ERROR_USERNAME_USER    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USESR_NOT_EXIST  = 1003
	ERROR_TOKEN_EXIST      = 1004
	ERROR_TOKEN_RUNTIME    = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_TYPE_WRONG = 1007
	ERROR_USER_NO_RIGHT    = 1008
	// code = 2000... 文章模块的错误
	ERROR_ART_NOT_EXIST = 2001

	// code = 3000... 分类模块的错误
	ERROR_CATENAME_USER  = 3001
	ERROR_CATE_NOT_EXIST = 3002

	// code = 4000... 权限模块的错误
	ERROR_ADD_POLICY          = 4001
	ERROR_ADD_GROUPING_POLICY = 4002
)

var codeMsg = map[int]string{
	SUCCSE:                    "OK",
	ERROR:                     "FAIL",
	ERROR_USERNAME_USER:       "用户名已存在！",
	ERROR_PASSWORD_WRONG:      "密码错误",
	ERROR_USESR_NOT_EXIST:     "用户不存在",
	ERROR_TOKEN_EXIST:         "TOKEN不存在",
	ERROR_TOKEN_RUNTIME:       "TOKEN已过期",
	ERROR_TOKEN_WRONG:         "TOKEN不正确",
	ERROR_TOKEN_TYPE_WRONG:    "TOKEN格式错误",
	ERROR_CATENAME_USER:       "分类已存在！",
	ERROR_CATE_NOT_EXIST:      "分类不存在！",
	ERROR_ART_NOT_EXIST:       "文章不存在",
	ERROR_USER_NO_RIGHT:       "该用户无权限",
	ERROR_ADD_POLICY:          "添加用户策略失败",
	ERROR_ADD_GROUPING_POLICY: "添加用户角色策略失败",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
