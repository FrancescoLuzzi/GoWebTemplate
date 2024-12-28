package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/interfaces"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

type (
	UserLogin struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}
	UserSignup struct {
		Email     string `validate:"required,email"`
		Password  string `validate:"required,password"`
		FirstName string `validate:"required,max=50"`
		LastName  string `validate:"required,max=50"`
	}
	UserUpdate struct {
		Email           string `validate:"required,email"`
		Password        string `validate:"required,password"`
		PasswordConfirm string
		FirstName       string `validate:"required,max=50"`
		LastName        string `validate:"required,max=50"`
	}
)

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

func (h *AuthHandler) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", h.handleLogin)
	mux.HandleFunc("POST /signup", h.handleSignup)
	mux.HandleFunc("GET /refresh", h.handleRefreshJWT)
	return mux
}

func (h *AuthHandler) handleSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var credentials UserSignup

	if err = h.decoder.Decode(&credentials, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.validate.Struct(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	passwordHash, err := auth.HashPassword(credentials.Password, &auth.DefaultConf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := types.User{
		Password:  passwordHash,
		Email:     credentials.Email,
		FirstName: credentials.Email,
		LastName:  credentials.LastName,
	}
	uid, err := h.service.Signup(r.Context(), &user)
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
	user, err := h.store.GetByEmail(r.Context(), &credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	valid, err := auth.ValidatePassword(credentials.Password, user.Password)
	if err != nil {
		slog.Info("couldn't validate password", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !valid {
		http.Error(w, "password not valid", http.StatusBadRequest)
		return
	}
	authToken, authExp, err := auth.CreateJWT(user.Id, auth.AuthToken, &h.cfg.JWTConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, refreshExp, err := auth.CreateJWT(user.Id, auth.RefreshToken, &h.cfg.JWTConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthTokenCookie,
		Value:    refreshToken,
		Expires:  refreshExp,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": authToken,
		"exp":   authExp,
	})
}

func (h *AuthHandler) handleRefreshJWT(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetRefreshToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := auth.ValidateJWT(refreshToken, &h.cfg.JWTConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	slog.Info("claims", "uid", claims["userId"])

	userId, err := uuid.Parse(claims["userId"].(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	authToken, authExp, err := auth.CreateJWT(userId, auth.AuthToken, &h.cfg.JWTConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": authToken,
		"exp":   authExp,
	})
}
