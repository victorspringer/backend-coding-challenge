package context

import (
	"context"

	"github.com/google/uuid"
)

// ContextKey is a custom type for context keys.
type ContextKey string

var (
	CTX_REQUEST_ID    ContextKey = "requestID"
	CTX_USER_LEVEL    ContextKey = "userLevel"
	CTX_USER_USERNAME ContextKey = "userUsername"
)

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(CTX_REQUEST_ID).(string)
	if !ok {
		requestID = uuid.New().String()
	}
	return requestID
}

// GetUserLevel retrieves the user's access level from the context.
func GetUserLevel(ctx context.Context) string {
	if level, ok := ctx.Value(CTX_USER_LEVEL).(string); ok {
		return level
	}
	return "anonymous"
}

// GetUserUsername retrieves the user's username from the context.
func GetUserUsername(ctx context.Context) string {
	if username, ok := ctx.Value(CTX_USER_USERNAME).(string); ok {
		return username
	}
	return ""
}
