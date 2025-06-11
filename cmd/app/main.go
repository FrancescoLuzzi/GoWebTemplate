package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/FrancescoLuzzi/GoWebTemplate/app"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/db"
	"github.com/FrancescoLuzzi/GoWebTemplate/app/middlewares"
	"github.com/FrancescoLuzzi/GoWebTemplate/public"
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
	appMux := middlewares.LoggingMiddleware(app.InitializeRoutes(conf, db))
	mux.Handle("/", appMux)
	mux.Handle("/public/assets/", public.FixCompressedContentHeaders(http.StripPrefix("/public/assets/", http.FileServerFS(public.AssetFs()))))
	fmt.Printf("Starting server %s:%s\n", conf.ServerConfig.Host, conf.ServerConfig.Port)
	http.ListenAndServe(conf.ServeAddr(), mux)
}
