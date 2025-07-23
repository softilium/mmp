package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigData struct {
	DbConnectionString             string
	dbDialect                      string
	ListenAddr                     string
	Debug                          bool
	FrontEndFolder                 string
	AdminEmail                     string
	AdminPassword                  string
	DbConnectionStringForMigration string
	ImagesFolder                   string
	ImagesFolderForMigration       string
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
	Cfg.DbConnectionStringForMigration = os.Getenv("DB_CONNECTION_STRING_FOR_MIGRATION")

	Cfg.FrontEndFolder, _ = strings.CutSuffix(Cfg.FrontEndFolder, "/")

	Cfg.ImagesFolder = os.Getenv("IMAGES_FOLDER")
	Cfg.ImagesFolderForMigration = os.Getenv("IMAGES_FOLDER_FOR_MIGRATION")

}
