package auth

import "context"

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func WithUserRole(ctx context.Context, userID, role string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, roleKey, role)
	return ctx
}

func GetUserFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(userIDKey)
	userId, ok := val.(string)
	return userId, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(roleKey)
	role, ok := val.(string)
	return role, ok
}
