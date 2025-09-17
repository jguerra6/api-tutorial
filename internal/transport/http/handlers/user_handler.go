package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jguerra6/api-tutorial/internal/platform/ctxutils"
	"github.com/jguerra6/api-tutorial/internal/transport/http/dto"
	"github.com/jguerra6/api-tutorial/internal/transport/http/writer"
)

type UsersService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, uid string) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, uid string) error
}

type UsersHandler struct {
	Svc UsersService
}

func NewUserHandler(svc UsersService) *UsersHandler {
	return &UsersHandler{Svc: svc}
}

func (h *UsersHandler) PublicCreateUser(w http.ResponseWriter, r *http.Request) {

	var req = &dto.CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writer.JSON(w, http.StatusBadRequest, map[string]any{
			"error":   "invalid_json",
			"message": "malformed JSON body",
		})
		return
	}

	req.Role = ""
	res, err := h.Svc.CreateUser(r.Context(), req)
	if err != nil {
		writer.WriteAppError(w, r, err)
		return
	}
	writer.JSON(w, http.StatusCreated, res)
}

func (h *UsersHandler) AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var req = &dto.CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writer.JSON(w, http.StatusBadRequest, map[string]any{
			"error":   "invalid_json",
			"message": "malformed JSON body",
		})
		return
	}

	res, err := h.Svc.CreateUser(r.Context(), req)
	if err != nil {
		writer.WriteAppError(w, r, err)
		return
	}
	writer.JSON(w, http.StatusCreated, res)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := ctxutils.Params(ctx)
	userId, ok := params["userId"]
	if !ok {
		writer.WriteAppError(w, r, fmt.Errorf("invalid userId"))
		return
	}

	err := h.Svc.DeleteUser(ctx, userId)
	if err != nil {
		writer.WriteAppError(w, r, err)
		return
	}

	writer.JSON(w, http.StatusNoContent, "")
}
