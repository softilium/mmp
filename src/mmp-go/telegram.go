package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var Bot *tgbotapi.BotAPI

func TelegramUserName(userId, userName string) string {
	if strings.TrimSpace(userName) == "" {
		return "tg." + userId
	} else {
		return userName
	}
}

func TelegramMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if after, ok := strings.CutPrefix(auth, "tg "); ok {
			data := after
			tgTokens := strings.Split(data, "~~")
			if len(tgTokens) != 3 {
				fmt.Printf("Invalid Telegram token format: %s\n", data)
				next.ServeHTTP(w, r)
				return
			}
			expIn := 24 * time.Hour
			err := initdata.Validate(tgTokens[0], Cfg.TgbotApiKey, expIn)
			if err != nil {
				fmt.Printf("Telegram authorization failed for token: %s\n", tgTokens[0])
				next.ServeHTTP(w, r)
				return
			}
			userId := tgTokens[1] // This is the Telegram user ID. We use it when user didn't set username in Telegram
			userName := tgTokens[2]
			users, _, err := models.Dbc.UserDef.SelectEntities([]*elorm.Filter{
				elorm.AddFilterEQ(models.Dbc.UserDef.TelegramVerified, true),
				elorm.AddFilterEQ(models.Dbc.UserDef.TelegramUsername, TelegramUserName(userId, userName)),
			}, nil, 1, 1)
			if err != nil {
				fmt.Printf("Error fetching user by Telegram username: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}
			if len(users) == 1 {
				ctx = context.WithValue(ctx, models.TelegramUserCtxKey, users[0])
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else if len(users) == 0 {
				fmt.Printf("No user found for Telegram username: %s\n", TelegramUserName(userId, userName))
				next.ServeHTTP(w, r)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func init() {
	if Cfg.TgbotApiKey == "" {
		return
	}

	var err error
	Bot, err = tgbotapi.NewBotAPI(Cfg.TgbotApiKey)
	if err != nil {
		log.Fatalf("Failed to connect to Telegram bot: %v", err)
	}

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 5

		updates := Bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatalf("Failed to get updates from Telegram bot: %v", err)
		}

		for update := range updates {
			if update.Message == nil || update.Message.From.IsBot {
				continue // ignore any non-Message updates as well as messages from bots
			}

			tgUserName := update.Message.From.UserName
			tgUserId := string(update.Message.From.ID)

			exUser, _, _err := models.Dbc.UserDef.SelectEntities([]*elorm.Filter{
				elorm.AddFilterEQ(models.Dbc.UserDef.TelegramUsername, TelegramUserName(tgUserId, tgUserName)),
				elorm.AddFilterEQ(models.Dbc.UserDef.TelegramVerified, true),
				elorm.AddFilterEQ(models.Dbc.UserDef.IsDeleted, false),
			}, nil, 1, 1)

			if _err != nil {
				log.Printf("Error fetching user by Telegram username: %v", _err)
				continue
			}
			if len(exUser) == 1 {
				if !exUser[0].IsActive() {
					continue
				}
				if exUser[0].TelegramChatId() != update.Message.Chat.ID {
					exUser[0].SetTelegramChatId(update.Message.Chat.ID)
					err = exUser[0].Save(context.Background())
					if err != nil {
						log.Printf("Error saving user chat ID: %v", err)
					} else {
						log.Printf("Updated chat ID for user %s: %d", exUser[0].Username(), update.Message.Chat.ID)
					}
				} else {
					newUser, err := models.Dbc.CreateUser()
					if err != nil {
						log.Printf("Error creating new user: %v", err)
						continue
					}
					newUser.SetTelegramUsername(TelegramUserName(tgUserId, tgUserName))
					newUser.SetUsername(TelegramUserName(tgUserId, tgUserName))
					newUser.SetTelegramVerified(true)
					newUser.SetEmail(TelegramUserName(tgUserId, tgUserName) + "@telegram.tg")
					newUser.SetTelegramChatId(update.Message.Chat.ID)
					newUser.SetPassword(elorm.NewRef())
					newUser.SetTelegramChatId(update.Message.Chat.ID)
				}

			}

			{ // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				Bot.Send(msg)
			}
		}
	}()
}
