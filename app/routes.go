package app

import (
	"net/http"

	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/middleware/htmx"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/services/user"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/utils"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/views/landing"
	"github.com/jmoiron/sqlx"
)

func InitializeRoutes(conf config.AppConfig, db *sqlx.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /home", htmx.TrapHxRequest(utils.RenderComponentHandler(landing.Index())))
	mux.Handle("GET /signup", htmx.TrapHxRequest(utils.RenderComponentHandler(landing.Signup())))
	mux.Handle("GET /login", htmx.TrapHxRequest(utils.RenderComponentHandler(landing.Login())))

	userHandler := user.NewHandler(user.NewUserStore(db))
	mux.Handle("/", userHandler.GetRoutes(conf))
	return mux
}
