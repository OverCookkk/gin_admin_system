package middleware

import (
	"gin_admin_system/internal/app/ginx"
	"gin_admin_system/pkg/auth"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

// func wrapUserAuthContext(c *gin.Context, userID uint64, userName string) {
//     ctx := contextx.NewUserID(c.Request.Context(), userID)
//     ctx = contextx.NewUserName(ctx, userName)
//     ctx = logger.NewUserIDContext(ctx, userID)
//     ctx = logger.NewUserNameContext(ctx, userName)
//     c.Request = c.Request.WithContext(ctx)
// }

func UserAuthMiddleware(a auth.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ginx.GetToken(c)
		claims, err := a.ParseToken(token)
		if err != nil {
			c.Abort()
			return
		}

		// 判断token是否过期
		if time.Now().Unix() > claims.ExpiresAt {
			// todo: token过期
			// app.ReturnWithDetailed()
			c.Abort()
		}

		// 获取到userID封装进其他组件中
		// tokenUserID = userID-username
		tokenUserID, err := a.ParseUserID(c.Request.Context(), token)
		if err != nil {
			c.Abort()
			return
		}
		idx := strings.Index(tokenUserID, "-")
		if idx == -1 {
			c.Abort()
			return
		}

		// todo:获取到userID封装进其他组件中
		_, _ = strconv.ParseUint(tokenUserID[:idx], 10, 64)
		// wrapUserAuthContext(c, userID, tokenUserID[idx+1:])
		c.Next()
	}
}
