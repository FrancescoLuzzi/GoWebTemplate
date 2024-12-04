package main

import (
	"log"

	"github.com/FrancescoLuzzi/AQuickQuestion/app"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/db"
	"github.com/FrancescoLuzzi/AQuickQuestion/public"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conf := config.Config()
	db, err := db.Open(conf.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	fiberApp := fiber.New(fiber.Config{
		CompressedFileSuffixes: map[string]string{
			"gzip": ".gz",
			"br":   ".br",
			"zstd": ".zst",
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(c.Response().StatusCode()).JSON(map[string]string{
				"message": err.Error(),
			})
		},
	})
	public.RegisterAssets(fiberApp)
	app.InitializeRoutes(fiberApp.Group("/"), conf, db)
	fiberApp.Listen(":8080")
}
