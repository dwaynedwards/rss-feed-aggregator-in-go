package context

import (
	"context"
	"net/http"
)

type contextKey int

const (
	userIDContextKey = contextKey(iota + 1)
)

func SetUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func SetUserIDToRequestContext(r *http.Request, userID int64) *http.Request {
	ctx := SetUserIDToContext(r.Context(), userID)
	return r.WithContext(ctx)
}

func UserIDFromContext(ctx context.Context) int64 {
	user, ok := ctx.Value(userIDContextKey).(int64)
	if !ok {
		return 0
	}
	return user
}
