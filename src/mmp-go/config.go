package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigData struct {
	DbConnectionString string
	dbDialect          string
	ListenAddr         string
	Debug              bool
	FrontEndFolder     string
	AdminEmail         string // Admin email for the default admin user (only for empty database)
	AdminPassword      string // Admin password for the default admin user (only for empty database)
	ImagesFolder       string
	TgbotApiKey        string
	CORSAlloweedHosts  []string // Allowed CORS hosts
}

var Cfg ConfigData

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	Cfg.DbConnectionString = os.Getenv("DB_CONNECTION_STRING")
	Cfg.dbDialect = os.Getenv("DB_DIALECT")
	Cfg.ListenAddr = os.Getenv("LISTEN_ADDR")
	Cfg.Debug = os.Getenv("DEBUG") == "true"
	Cfg.FrontEndFolder = os.Getenv("FRONT_END_FOLDER")
	Cfg.AdminEmail = os.Getenv("ADMIN_EMAIL")
	Cfg.AdminPassword = os.Getenv("ADMIN_PASSWORD")

	Cfg.FrontEndFolder, _ = strings.CutSuffix(Cfg.FrontEndFolder, "/")

	Cfg.ImagesFolder = os.Getenv("IMAGES_FOLDER")
	Cfg.TgbotApiKey = os.Getenv("TGBOT_APIKEY")

	cah := os.Getenv("CORS_ALLOWED_HOSTS")
	Cfg.CORSAlloweedHosts = strings.Split(cah, ",")

}
