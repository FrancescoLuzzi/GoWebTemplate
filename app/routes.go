package app

import (
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/middleware/htmx"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/services/user"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/utils"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/views/landing"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func InitializeRoutes(group fiber.Router, conf config.AppConfig, db *sqlx.DB) {
	group.Use(htmx.New())
	group.Get("/", utils.RenderComponentHandler(landing.Index()))
	group.Get("/signup", utils.RenderComponentHandler(landing.Signup()))
	group.Get("/login", utils.RenderComponentHandler(landing.Login()))

	userHandler := user.NewHandler(user.NewUserStore(db))
	userHandler.RegisterRoutes(group, conf)
}
