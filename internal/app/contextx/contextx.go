package contextx

import "context"

const (
	contextUserID = "user_id"
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
