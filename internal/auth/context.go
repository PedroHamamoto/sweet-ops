package auth

import "context"

type ctxKey string

const (
	ctxKeyUserID ctxKey = "user_id"
)

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ctxKeyUserID).(string)
	return userID, ok
}
