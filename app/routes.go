package app

import (
	"net/http"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/handlers"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/middlewares"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/services"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/stores"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/utils"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/views/landing"
	"github.com/jmoiron/sqlx"
)

func InitializeRoutes(conf config.AppConfig, db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()

	userStore := stores.NewUserStore(db)

	authService := services.NewAuthService(userStore, &conf)
	userService := services.NewUserService(userStore)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	authMiddleware := middlewares.NewAuthMiddleware(userStore, &conf.JWTConfig)

	md := middlewares.Combine(
		middlewares.HxRequestMiddleware,
		authMiddleware,
	)
	mux.Handle("GET /", md(utils.RenderComponentHandler(landing.Index())))
	mux.Handle("GET /home", md(utils.RenderComponentHandler(landing.Index())))
	mux.Handle("GET /signup", md(utils.RenderComponentHandler(landing.Signup())))
	mux.Handle("GET /login", md(utils.RenderComponentHandler(landing.Login())))

	authHandler.RegisterRoutes(mux, md)
	userHandler.RegisterRoutes(mux, md)
	return mux
}
