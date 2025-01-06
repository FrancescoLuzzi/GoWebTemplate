package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type UserHandler struct {
	service  interfaces.UserService
	validate *validator.Validate
	decoder  *schema.Decoder
}

func NewUserHandler(service interfaces.UserService) UserHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	registerPasswordValidation(validate)
	return UserHandler{
		service:  service,
		validate: validate,
		decoder:  schema.NewDecoder(),
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware middlewares.Middleware) {
	// admin routes
	mux.Handle("GET /user/profile", authMiddleware(http.HandlerFunc(h.handleCurrentUserProfile)))
}

func (h *UserHandler) handleCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := auth.UserFromCtx(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := h.service.GetById(ctx, &userId)
	if err != nil {
		slog.Info("failed to get user by id", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
