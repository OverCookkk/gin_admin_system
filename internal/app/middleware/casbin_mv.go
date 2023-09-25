package middleware

import (
	"gin_admin_system/internal/app/contextx"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		m := c.Request.Method
		userID := contextx.GetUserID(c.Request.Context())
		success, err := enforcer.Enforce(strconv.FormatUint(userID, 10), p, m)
		if err != nil {
			return
		}
		if !success { // 没有权限
			// todo:权限不足
			// app.ReturnWithDetailed()
			return
		}
		c.Next()
	}
}
