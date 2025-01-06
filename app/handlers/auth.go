package handlers

import (
	"encoding/json"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type UserLogin struct {
	Email    string `json:"email" schema:"email" validate:"required,email"`
	Password string `json:"password" schema:"password" validate:"required"`
}

func (u UserLogin) ToUser() *types.User {
	return &types.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserInfos struct {
	Email     string `json:"email" schema:"email" validate:"required,email"`
	Password  string `json:"password" schema:"password" validate:"required,password"`
	FirstName string `json:"first_name" schema:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" schema:"last_name" validate:"required,max=50"`
}

func (u UserInfos) ToUser() *types.User {
	return &types.User{
		Password:  u.Password,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

}

func registerPasswordValidation(v *validator.Validate) {
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		var (
			hasNumber      = false
			hasSpecialChar = false
			hasLetter      = false
			hasSuitableLen = false
		)

		password := fl.Field().String()

		if utf8.RuneCountInString(password) <= 30 && utf8.RuneCountInString(password) >= 8 {
			hasSuitableLen = true
		}

		for _, c := range password {
			switch {
			case unicode.IsNumber(c):
				hasNumber = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSpecialChar = true
			case unicode.IsLetter(c) || c == ' ':
				hasLetter = true
			default:
				return false
			}
		}

		return hasNumber && hasSpecialChar && hasLetter && hasSuitableLen
	})
}

type AuthHandler struct {
	service  interfaces.AuthService
	validate *validator.Validate
	decoder  *schema.Decoder
	cfg      *config.AppConfig
}

func NewAuthHandler(service interfaces.AuthService) AuthHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	registerPasswordValidation(validate)
	return AuthHandler{
		service:  service,
		validate: validate,
		decoder:  schema.NewDecoder(),
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.handleLogin)
	mux.HandleFunc("POST /auth/signup", h.handleSignup)
	mux.HandleFunc("GET /auth/refresh", h.handleRefreshJWT)
}

func (h *AuthHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var credentials UserInfos

	if err = h.decoder.Decode(&credentials, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.validate.StructCtx(r.Context(), &credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uid, err := h.service.Signup(r.Context(), credentials.ToUser())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"userId": uid.String(),
	})
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var credentials UserLogin
	if err = h.decoder.Decode(&credentials, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.service.Login(r.Context(), credentials.ToUser())
	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenCookie,
		Value:    res.RefreshToken.Token,
		Expires:  res.RefreshToken.Exp,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": res.AuthToken.Token,
		"exp":   res.AuthToken.Exp,
	})
}

func (h *AuthHandler) handleRefreshJWT(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetRefreshToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	authToken, err := h.service.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": authToken.Token,
		"exp":   authToken.Exp,
	})
}
