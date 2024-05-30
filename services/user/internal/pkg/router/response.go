package router

import (
	"encoding/json"
	"net/http"

	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/log"
)

type response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"response"`
}

func (rt *router) respond(w http.ResponseWriter, r *http.Request, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")

	res := response{
		StatusCode: code,
		Body:       body,
	}

	b, err := json.Marshal(res)
	if err != nil {
		rt.logger.Error("failed to unmarshal response", log.Error(err), log.String("requestId", getRequestID(r.Context())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(b)
}
