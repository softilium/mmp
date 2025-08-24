package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	EMAIL_MODE_NONE   = 0
	EMAIL_MODE_DIRECT = 1
	EMAIL_MODE_RESEND = 2
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

	EmailMode int // 1 = direct SMTP, 2 = Resend

	EmailSMTPUserName string
	EmailSMTPPassword string
	EmailSMTPHost     string
	EmailSMTPPort     int
	EmailFrom         string

	Email_Resend_ApiKey string
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

	Cfg.EmailSMTPUserName = os.Getenv("EMAIL_SMTP_USERNAME")
	Cfg.EmailSMTPPassword = os.Getenv("EMAIL_SMTP_PASSWORD")
	Cfg.EmailSMTPHost = os.Getenv("EMAIL_SMTP_HOST")
	Cfg.EmailSMTPPort = 587
	if portStr := os.Getenv("EMAIL_SMTP_PORT"); portStr != "" {
		_, _ = fmt.Sscanf(portStr, "%d", &Cfg.EmailSMTPPort)
	}
	Cfg.EmailFrom = os.Getenv("EMAIL_FROM")

	Cfg.Email_Resend_ApiKey = os.Getenv("EMAIL_RESEND_APIKEY")

	emailMode := os.Getenv("EMAIL_MODE")
	Cfg.EmailMode = EMAIL_MODE_NONE
	if emailMode == "direct" {
		Cfg.EmailMode = EMAIL_MODE_DIRECT
	}
	if emailMode == "resend" {
		Cfg.EmailMode = EMAIL_MODE_RESEND
	}

}
