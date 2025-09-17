package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jguerra6/api-tutorial/internal/transport/http/writer"
)

type DBPinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	DB DBPinger
}

func NewHealthHandler(db DBPinger) *HealthHandler { return &HealthHandler{DB: db} }

func (h *HealthHandler) Healthz(w http.ResponseWriter, _ *http.Request) {
	writer.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		writer.JSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "degraded",
			"db":     "not_configured",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		writer.JSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "degraded",
			"db":     "unreachable",
			"error":  err.Error(),
		})
		return
	}

	writer.JSON(w, http.StatusOK, map[string]any{
		"status": "ready",
		"db":     "ok",
	})
}
