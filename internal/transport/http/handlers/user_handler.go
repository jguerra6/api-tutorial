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

// PublicCreateUser godoc
// @Summary Create a user (public)
// @Description Create a new user. Role will be blank/ignored for public creation.
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create User Request"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users [post]
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

// AdminCreateUser godoc
// @Summary Create a user (admin)
// @Description Create a new user as an administrator. Requires authentication.
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create User Request"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/users/admin [post]
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID. Requires authentication.
// @Tags Users
// @Produce json
// @Param userId path string true "User ID"
// @Success 204 {string} string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/users/{userId} [delete]
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
