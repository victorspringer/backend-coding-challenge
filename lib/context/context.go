package context

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

var (
	CTX_REQUEST_ID ContextKey = "requestID"
	CTX_USER_LEVEL ContextKey = "userLevel"
)

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(CTX_REQUEST_ID).(string)
	if !ok {
		requestID = uuid.New().String()
	}
	return requestID
}
