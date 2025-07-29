package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type TokenItem struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

var TokensByAT = make(map[string]TokenItem)

func UserFromHttpRequest(r *http.Request) (*User, error) {
	raw := r.Header.Get("Authorization")
	if raw == "" {
		return nil, fmt.Errorf("no auth header")
	}
	if len(raw) < 7 || raw[:7] != "Bearer " {
		return nil, fmt.Errorf("invalid auth header")
	}
	token := raw[7:]

	var foundToken *TokenItem
	var foundUserRef string
	for k, item := range TokensByAT {
		if item.AccessToken == token {
			if !item.AccessTokenExpiresAt.After(time.Now()) {
				return nil, fmt.Errorf("access token expired")
			}
			foundToken = &item
			foundUserRef = k
			break
		}
	}

	if foundToken != nil {
		user, err := Dbc.LoadUser(foundUserRef)
		if err != nil {
			return nil, fmt.Errorf("unable to load user: %v", err)
		}
		if !user.IsActive() {
			return nil, fmt.Errorf("user isn't active")
		}
		return user, nil
	}
	return nil, fmt.Errorf("unregistered token")
}

func GenerateToken(user *User) TokenItem {
	newToken := TokenItem{
		AccessToken:           uuid.NewString(),
		RefreshToken:          uuid.NewString(),
		AccessTokenExpiresAt:  time.Now().Add(20 * time.Second),
		RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour * 90),
	}
	TokensByAT[user.RefString()] = newToken
	return newToken
}
