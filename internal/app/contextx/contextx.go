package contextx

import "context"

type (
	userIDContext  struct{}
	traceIDContext struct{}
	transContext   struct{}

	// todo
	// transCtx     struct{}
	// noTransCtx   struct{}
	// transLockCtx struct{}
	// userIDCtx    struct{}
	// userNameCtx  struct{}
	// traceIDCtx   struct{}
)

func SetUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDContext{}, userID)
}

func GetUserID(ctx context.Context) uint64 {
	v := ctx.Value(userIDContext{})
	if v != nil {
		if s, ok := v.(uint64); ok {
			return s
		}
	}
	return 0
}

func SetTraceID(ctx context.Context, TraceID string) context.Context {
	return context.WithValue(ctx, traceIDContext{}, TraceID)
}

func GetTraceID(ctx context.Context) string {
	v := ctx.Value(traceIDContext{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// SetTrans 把事务对象存进上下文中
func SetTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transContext{}, trans)
}

func GetTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transContext{})
	return v, v != nil
}
