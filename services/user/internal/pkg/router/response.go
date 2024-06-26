package router

import (
	"encoding/json"
	"net/http"

	"github.com/victorspringer/backend-coding-challenge/lib/context"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
)

type response struct {
	StatusCode int         `json:"statusCode"`
	Response   interface{} `json:"response,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func (rt *router) respond(w http.ResponseWriter, r *http.Request, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")

	var res response

	if code >= 400 {
		res = response{
			StatusCode: code,
			Error:      body.(string),
		}
	} else {
		res = response{
			StatusCode: code,
			Response:   body,
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		rt.logger.Error("failed to unmarshal response", log.Error(err), log.String("requestId", context.GetRequestID(r.Context())))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(b)
}
