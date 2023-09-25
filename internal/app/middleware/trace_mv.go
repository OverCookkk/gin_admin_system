package middleware

import (
	"gin_admin_system/internal/app/contextx"
	"gin_admin_system/pkg/logger"
	"gin_admin_system/pkg/util/trace"
	"github.com/gin-gonic/gin"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = trace.NewTraceID()
		}

		ctx := contextx.SetTraceID(c.Request.Context(), traceID)
		ctx = logger.SetTraceIDContext(ctx, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-Id", traceID)

		c.Next()
	}
}
