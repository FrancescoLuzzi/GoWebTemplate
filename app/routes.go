package app

import (
	"net/http"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/cache"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/handlers"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/middlewares"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/services"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/stores"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/utils"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/views/landing"
	"github.com/jmoiron/sqlx"
)

func InitializeRoutes(conf config.AppConfig, cache cache.Cache, db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()

	userStore := stores.NewUserStore(db)

	authService := services.NewAuthService(userStore, &conf)
	userService := services.NewUserService(userStore)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	authMiddleware := middlewares.NewAuthMiddleware(userStore, &conf.JWTConfig)
	cacheMiddleware := middlewares.CacheInjectorMiddleware(cache)

	md := middlewares.Combine(
		cacheMiddleware,
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
