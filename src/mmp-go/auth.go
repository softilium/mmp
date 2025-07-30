package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/softilium/elorm"
	"golang.org/x/crypto/bcrypt"
)

type TelegramUserCtxKeyType string

const TelegramUserCtxKey TelegramUserCtxKeyType = "tguser"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func UserFromHttpRequest(r *http.Request) (*User, *Token, error) {
	tgUser := r.Context().Value(TelegramUserCtxKey)
	if tgUser != nil {
		if user, ok := tgUser.(*User); ok {
			if user.IsActive() {
				return user, nil, nil
			}
			return nil, nil, fmt.Errorf("user isn't active")
		}
		return nil, nil, fmt.Errorf("invalid user context type for tg")
	}
	raw := r.Header.Get("Authorization")
	if raw == "" {
		return nil, nil, fmt.Errorf("no auth header")
	}
	if len(raw) < 7 || raw[:7] != "Bearer " {
		return nil, nil, fmt.Errorf("invalid auth header")
	}
	token := raw[7:]
	dbtokens, _, err := DB.TokenDef.SelectEntities([]*elorm.Filter{elorm.AddFilterEQ(DB.TokenDef.AccessToken, token)}, nil, 1, 1)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load token from db: %v", err)
	}
	if len(dbtokens) == 1 {
		if dbtokens[0].AccessTokenExpiresAt().Before(time.Now()) {
			return nil, nil, fmt.Errorf("access token expired")
		}
		if !dbtokens[0].User().IsActive() {
			return nil, nil, fmt.Errorf("user isn't active")
		}
		return dbtokens[0].User(), dbtokens[0], nil
	}
	return nil, nil, fmt.Errorf("unregistered token")
}

func GenerateToken(user *User) (*Token, error) {
	newToken, err := DB.CreateToken()
	if err != nil {
		return nil, fmt.Errorf("unable to create token: %v", err)
	}
	newToken.SetAccessToken(elorm.NewRef())
	newToken.SetRefreshToken(elorm.NewRef())
	if Cfg.Debug {
		newToken.SetAccessTokenExpiresAt(time.Now().Add(20 * time.Minute))
	} else {
		newToken.SetAccessTokenExpiresAt(time.Now().Add(20 * time.Second))
	}
	newToken.SetRefreshTokenExpiresAt(time.Now().Add(24 * time.Hour * 90))
	newToken.SetUser(user)
	err = newToken.Save(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to save token: %v", err)
	}
	return newToken, nil
}
