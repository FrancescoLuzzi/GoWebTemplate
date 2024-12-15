package main

import (
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/FrancescoLuzzi/AQuickQuestion/app"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/db"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/middleware/logging"
	"github.com/FrancescoLuzzi/AQuickQuestion/public"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func setupGlobalLogger(w io.Writer, cfg *config.ServerConfig) {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})
	handler := slog.New(jsonHandler)
	slog.SetDefault(handler)
}

func main() {
	conf := config.Config()
	setupGlobalLogger(os.Stdout, &conf.ServerConfig)
	db, err := db.Open(conf.DbConfig)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	appMux := app.InitializeRoutes(conf, db)
	mux.Handle("/", appMux)
	mux.Handle("/public/assets/", public.FixCompressedContentHeaders(http.StripPrefix("/public/assets/", http.FileServerFS(public.AssetFs()))))

	http.ListenAndServe(":8080", logging.NewLoggingMiddleware(mux))
}
