package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/middlewares"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/utils"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/views/landing"
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

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, md middlewares.Middleware) {
	// admin routes
	mux.Handle("GET /user/profile", md(http.HandlerFunc(h.handleCurrentUserProfile)))
	mux.Handle("POST /user/profile", md(http.HandlerFunc(h.handleProfileUpdate)))
	mux.Handle("POST /user/password", md(http.HandlerFunc(h.handlePasswordUpdate)))
	mux.Handle("GET /profile", md(http.HandlerFunc(h.handleCurrentUserProfileView)))
}

func (h *UserHandler) currentUser(w http.ResponseWriter, r *http.Request) (types.User, error) {
	ctx := r.Context()
	userId, err := auth.UserFromCtx(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return types.User{}, err
	}
	return h.service.GetById(ctx, &userId)

}

func (h *UserHandler) handleCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(w, r)
	if err != nil {
		slog.Info("failed to get user by id", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleCurrentUserProfileView(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(w, r)
	if err != nil {
		slog.Info("failed to get user by id", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RenderComponentHandler(landing.Profile(&user)).ServeHTTP(w, r)
}

func (h *UserHandler) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// read post body
	// update username, email, ...
}

func (h *UserHandler) handlePasswordUpdate(w http.ResponseWriter, r *http.Request) {
	// read post body
	// check if old password is correct
	// hash and save new password
}
