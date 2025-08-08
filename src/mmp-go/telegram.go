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
			users, _, err := DB.UserDef.SelectEntities([]*elorm.Filter{
				elorm.AddFilterEQ(DB.UserDef.TelegramVerified, true),
				elorm.AddFilterEQ(DB.UserDef.TelegramUsername, TelegramUserName(userId, userName)),
			}, nil, 1, 1)
			if err != nil {
				fmt.Printf("Error fetching user by Telegram username: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}
			if len(users) == 1 {
				ctx = context.WithValue(ctx, TelegramUserCtxKey, users[0])
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

	// bot mesages handler is responsible for 3 things:
	// 1. Updating user chat ID in the database if it has changed
	// 2. Creating a new user if it doesn't exist
	// 3. Forwarding messages to the bot to admins
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
			tgUserId := fmt.Sprintf("%d", update.Message.From.ID)
			exUsers, _, _err := DB.UserDef.SelectEntities([]*elorm.Filter{
				elorm.AddFilterEQ(DB.UserDef.TelegramUsername, TelegramUserName(tgUserId, tgUserName)),
				elorm.AddFilterEQ(DB.UserDef.TelegramVerified, true),
			}, nil, 1, 1)

			if _err != nil {
				log.Printf("Error fetching user by Telegram username: %v", _err)
				continue
			}
			if len(exUsers) == 1 {
				if !exUsers[0].IsActive() {
					continue
				}
				if exUsers[0].TelegramChatId() != update.Message.Chat.ID {
					exUsers[0].SetTelegramChatId(update.Message.Chat.ID)
					err = exUsers[0].Save(context.Background())
					if err != nil {
						log.Printf("Error saving user chat ID: %v", err)
					} else {
						log.Printf("Updated chat ID for user %s: %d", exUsers[0].Username(), update.Message.Chat.ID)
					}
				}

			} else {
				newUser, err := DB.CreateUser()
				if err != nil {
					log.Printf("Error creating new user: %v", err)
					continue
				}
				newPassword := elorm.NewRef()
				newPasswordHash, err := HashPassword(newPassword)
				if err != nil {
					log.Printf("Error hashing password: %v", err)
					continue
				}
				newUser.SetTelegramUsername(TelegramUserName(tgUserId, tgUserName))
				newUser.SetUsername(TelegramUserName(tgUserId, tgUserName))
				newUser.SetTelegramVerified(true)
				newUser.SetEmail(TelegramUserName(tgUserId, tgUserName) + "@telegram.tg")
				newUser.SetTelegramChatId(update.Message.Chat.ID)
				newUser.SetPasswordHash(newPasswordHash)
				newUser.SetTelegramChatId(update.Message.Chat.ID)
				err = newUser.Save(context.Background())
				if err != nil {
					log.Printf("Error saving new user: %v", err)
					continue
				}

				_, _ = Bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					fmt.Sprintf(`

Вас приветствует бот сервиса RiverStores. Вы можете использовать сервис, нажав на название бота сверху и открыв мини-приложение по ссылке. Либо, в списке чатов или в окне чата нажать кнопку (ОТКРЫТЬ/OPEN).

Вы можете открыть ссылку в браузере: https://rives-stores.com
Ваш логин  : %s
Ваш пароль : %s

Есть вопросы, проблемы, предложения, идеи, отзывы? Просто напишите их в чат боту и они будут сразу переданы администратору сервиса.
`,
						newUser.Username(), newPassword)))

			}

			admins, _, err := DB.UserDef.SelectEntities([]*elorm.Filter{
				elorm.AddFilterEQ(DB.UserDef.TelegramVerified, true),
				elorm.AddFilterNOEQ(DB.UserDef.TelegramUsername, ""),
				elorm.AddFilterEQ(DB.UserDef.IsActive, true),
				elorm.AddFilterEQ(DB.UserDef.Admin, true),
			}, nil, 0, 0)
			if err != nil {
				log.Printf("Error fetching admin users: %v", err)
				continue
			}
			for _, admin := range admins {
				msg := tgbotapi.NewMessage(admin.TelegramChatId(), fmt.Sprintf("Message from %s:\n\r%s", update.Message.From.UserName, update.Message.Text))
				msg.ReplyToMessageID = update.Message.MessageID
				_, err = Bot.Send(msg)
				if err != nil {
					log.Printf("Error sending message to admin %s: %v", admin.Username(), err)
				}
			}
		}

	}()
}
