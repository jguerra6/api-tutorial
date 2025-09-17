package writer

import (
	"encoding/json"
	"net/http"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
)

type errorEnvelope struct {
	Error     string            `json:"error"`
	Code      string            `json:"code,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, errorEnvelope{Error: msg})
}

func requestID(r *http.Request) string {
	if rid := ctxutils.RequestID(r.Context()); rid != "" {
		return rid
	}
	if rid := r.Header.Get("X-Request-ID"); rid != "" {
		return rid
	}
	return ""
}
