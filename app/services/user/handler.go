package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/services/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/go-playground/validator/v10"
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

var validate = validator.New(validator.WithRequiredStructEnabled())
var decoder = schema.NewDecoder()

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

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) Handler {
	return Handler{store}
}

func (h *Handler) GetRoutes(cfg config.AppConfig) *http.ServeMux {
	mux := http.NewServeMux()
	registerPasswordValidation(validate)
	mux.HandleFunc("POST /login", h.handleLogin(&cfg.JWTConfig))
	mux.HandleFunc("POST /signup", h.handleSignup)

	withJWT := auth.CreateJWTAuthHandler(h.store, &cfg.JWTConfig)
	// admin routes
	mux.HandleFunc("GET /profile", withJWT(h.handleCurrentUserProfile))
	return mux
}

func (h *Handler) handleSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var credentials UserSignup

	if err = decoder.Decode(&credentials, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = validate.Struct(&credentials); err != nil {
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
	uid, err := h.store.Create(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"userId": uid.String(),
	})
}

func (h *Handler) handleLogin(cfg *config.JWTConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var credentials UserLogin
		if err = decoder.Decode(&credentials, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(&credentials); err != nil {
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
		authToken, authExp, err := auth.CreateJWT(user.Id, auth.AuthToken, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		refreshToken, refreshExp, err := auth.CreateJWT(user.Id, auth.RefreshToken, cfg)
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
}

func (h *Handler) handleCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := auth.UserFromCtx(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := h.store.GetById(ctx, &userId)
	if err != nil {
		slog.Info("failed to get user by id", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
