package user

import (
	"log"
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/services/auth"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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

func (h *Handler) RegisterRoutes(router fiber.Router, cfg config.AppConfig) {
	registerPasswordValidation(validate)
	router.Post("/login", h.handleLogin(&cfg.JWTConfig))
	router.Post("/signup", h.handleSignup)

	// admin routes
	router.Get("/profile", h.handleCurrentUserProfile, auth.WithJWTAuth(h.store, &cfg.JWTConfig))
}

func (h *Handler) handleSignup(ctx fiber.Ctx) error {
	var credentials UserSignup
	if err := ctx.Bind().Body(&credentials); err != nil {
		return err
	}
	if err := validate.Struct(&credentials); err != nil {
		return err
	}
	passwordHash, err := auth.HashPassword(credentials.Password, &auth.DefaultConf)
	if err != nil {
		return err
	}
	user := types.User{
		Password:  passwordHash,
		Email:     credentials.Email,
		FirstName: credentials.Email,
		LastName:  credentials.LastName,
	}
	uid, err := h.store.Create(ctx.Context(), &user)
	if err != nil {
		return err
	}
	return ctx.JSON(map[string]string{
		"userId": uid.String(),
	})
}

func (h *Handler) handleLogin(cfg *config.JWTConfig) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var credentials UserLogin
		if err := ctx.Bind().Body(&credentials); err != nil {
			return err
		}
		if err := validate.Struct(&credentials); err != nil {
			return err
		}
		user, err := h.store.GetByEmail(ctx.Context(), &credentials.Email)
		if err != nil {
			return err
		}
		valid, err := auth.ValidatePassword(credentials.Password, user.Password)
		if err != nil {
			log.Printf("%v\n", err)
			ctx.Status(http.StatusBadRequest)
			return err
		}
		if !valid {
			log.Println("password not valid")
			ctx.Status(http.StatusBadRequest)
			return nil
		}
		authToken, authExp, err := auth.CreateJWT(user.Id, auth.AuthToken, cfg)
		if err != nil {
			return err
		}
		refreshToken, refreshExp, err := auth.CreateJWT(user.Id, auth.RefreshToken, cfg)
		if err != nil {
			return err
		}
		ctx.Cookie(&fiber.Cookie{
			Name:     auth.AuthTokenCookie,
			Value:    refreshToken,
			Expires:  refreshExp,
			SameSite: "strict",
			HTTPOnly: true,
			Secure:   true,
		})
		return ctx.JSON(map[string]any{
			"token": authToken,
			"exp":   authExp,
		})
	}
}

func (h *Handler) handleCurrentUserProfile(ctx fiber.Ctx) error {
	userId, err := auth.UserFromCtx(ctx)
	if err != nil {
		return err
	}
	user, err := h.store.GetById(ctx.Context(), &userId)
	if err != nil {
		log.Printf("failed to get user by id: %v\n", err)
		ctx.Status(http.StatusBadRequest)
		return err
	}
	return ctx.JSON(user)
}
