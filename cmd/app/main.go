package main

import (
	"log"
	"net/http"

	"github.com/FrancescoLuzzi/AQuickQuestion/app"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/config"
	"github.com/FrancescoLuzzi/AQuickQuestion/app/db"
	"github.com/FrancescoLuzzi/AQuickQuestion/public"
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
	mux := http.NewServeMux()
	appMux := app.InitializeRoutes(conf, db)
	mux.Handle("/", appMux)
	mux.Handle("/public/assets/", public.FixCompressedContentHeaders(http.StripPrefix("/public/assets/", http.FileServerFS(public.AssetFs()))))
	http.ListenAndServe(":8080", mux)
}
