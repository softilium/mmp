package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func UserFromHttpRequest(r *http.Request) (*User, error) {
	raw := r.Header.Get("Authorization")
	if raw == "" {
		return nil, fmt.Errorf("Unauthorized")
	}
	if len(raw) < 7 || raw[:7] != "Bearer " {
		return nil, fmt.Errorf("Unauthorized")
	}
	token := raw[7:]

	if item, ok := TokensByAT[token]; ok {
		if !item.AccessTokenExpiresAt.After(time.Now()) {
			return nil, fmt.Errorf("access token expired")
		}
		user, err := Dbc.LoadUser(item.UserRef)
		if err != nil {
			return nil, fmt.Errorf("Unauthorized")
		}
		if !user.IsActive() {
			return nil, fmt.Errorf("user isn't active")
		}
		return user, nil
	}
	return nil, fmt.Errorf("Unauthorized")
}

type TokenItem struct {
	UserRef               string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

var TokensByAT = make(map[string]TokenItem)

func GenerateToken(user *User) (TokenItem, string) {
	newToken := TokenItem{
		UserRef:               user.RefString(),
		RefreshToken:          uuid.NewString(),
		AccessTokenExpiresAt:  time.Now().Add(1 * time.Hour),
		RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour * 30),
	}

	res := uuid.NewString()

	TokensByAT[res] = newToken
	return newToken, res
}
