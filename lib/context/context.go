package context

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

var (
	CTX_REQUEST_ID    ContextKey = "requestID"
	CTX_USER_LEVEL    ContextKey = "userLevel"
	CTX_USER_USERNAME ContextKey = "userUsername"
)

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(CTX_REQUEST_ID).(string)
	if !ok {
		requestID = uuid.New().String()
	}
	return requestID
}

func GetUserLevel(ctx context.Context) string {
	if level, ok := ctx.Value(CTX_USER_LEVEL).(string); ok {
		return level
	}
	return "anonymous"
}

func GetUserUsername(ctx context.Context) string {
	if username, ok := ctx.Value(CTX_USER_USERNAME).(string); ok {
		return username
	}
	return ""
}
