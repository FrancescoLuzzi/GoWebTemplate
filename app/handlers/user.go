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

type UserUpdateInfos struct {
	Email     string `json:"email" schema:"email" validate:"required,email"`
	FirstName string `json:"first_name" schema:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" schema:"last_name" validate:"required,max=50"`
}

func (u *UserUpdateInfos) updateUser(user *types.User) *types.User {
	user.Email = u.Email
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	return user
}

type UserUpdatePassword struct {
	OldPassword string `json:"old_password" schema:"old_password" validate:"password"`
	NewPassword string `json:"new_password" schema:"new_password" validate:"password"`
}

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
	mdMustAuth := middlewares.Combine(
		middlewares.MustAuthMiddleware,
		md,
	)
	mux.Handle("GET /user/profile", mdMustAuth(http.HandlerFunc(h.handleCurrentUserProfile)))
	mux.Handle("POST /user/profile", mdMustAuth(http.HandlerFunc(h.handleProfileUpdate)))
	mux.Handle("POST /user/password", mdMustAuth(http.HandlerFunc(h.handlePasswordUpdate)))
	mux.Handle("GET /profile", md(http.HandlerFunc(h.handleCurrentUserProfileView)))
}

func (h *UserHandler) currentUser(w http.ResponseWriter, r *http.Request) (types.User, error) {
	ctx := r.Context()
	userId, err := auth.UserFromCtx(ctx)
	if err != nil {
		slog.Info("failed to get user by id", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return types.User{}, err
	}
	return h.service.GetById(ctx, &userId)

}

func (h *UserHandler) handleCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var infos UserUpdateInfos
	if err = h.decoder.Decode(&infos, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.validate.StructCtx(r.Context(), &infos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.currentUser(w, r)
	if err != nil {
		slog.Info("failed to get user by id", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.Update(r.Context(), infos.updateUser(&user))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) handlePasswordUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var passwords UserUpdatePassword
	if err = h.decoder.Decode(&passwords, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.validate.StructCtx(r.Context(), &passwords); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.currentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.UpdatePassword(r.Context(), &user, passwords.OldPassword, passwords.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) handleCurrentUserProfileView(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RenderComponentHandler(landing.Profile(&user)).ServeHTTP(w, r)
}
