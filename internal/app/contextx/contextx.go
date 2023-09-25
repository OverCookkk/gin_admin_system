package contextx

import "context"

const (
	contextUserID  = "user_id"
	contextTraceID = "trace_id"
)

func SetUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, contextUserID, userID)
}

func GetUserID(ctx context.Context) uint64 {
	v := ctx.Value(contextUserID)
	if v != nil {
		if s, ok := v.(uint64); ok {
			return s
		}
	}
	return 0
}

func SetTraceID(ctx context.Context, TraceID string) context.Context {
	return context.WithValue(ctx, contextTraceID, TraceID)
}

func GetTraceID(ctx context.Context) string {
	v := ctx.Value(contextTraceID)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
