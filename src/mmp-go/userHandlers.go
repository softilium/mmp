package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
)

type userPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	payload := userPayLoad{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	exist, _, err := dbc.UserDef.SelectEntities(
		[]*elorm.Filter{elorm.AddFilterEQ(dbc.UserDef.Email, payload.Email)}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(exist) > 0 {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("user with this email already exists: %s", payload.Email))
		return
	}
	newUser, err := dbc.CreateUser()
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	newUser.SetUsername(payload.Email)
	newUser.SetEmail(payload.Email)
	newUser.SetPassword(payload.Password)
	err = newUser.Save()
	if err != nil {
		HandleErr(w, 0, err)
		return
	}

	_, _ = generateToken(newUser)

	w.WriteHeader(http.StatusCreated)

}

type tokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	payload := userPayLoad{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	users, _, err := dbc.UserDef.SelectEntities(
		[]*elorm.Filter{elorm.AddFilterEQ(dbc.UserDef.Email, payload.Email)}, nil, 1, 1)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(users) != 1 {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	if !users[0].IsActive() {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("user is not active"))
		return
	}

	newToken, at := generateToken(users[0])
	w.WriteHeader(http.StatusOK)
	res := tokensResponse{AccessToken: at, RefreshToken: newToken.refreshToken}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	user, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}
	for token, item := range tokensByAT {
		if item.userRef == user.RefString() {
			delete(tokensByAT, token)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func generateToken(user *models.User) (TokenItem, string) {
	newToken := TokenItem{
		userRef:               user.RefString(),
		refreshToken:          uuid.NewString(),
		accessTokenExpiresAt:  time.Now().Add(5 * time.Minute),
		refreshTokenExpiresAt: time.Now().Add(24 * time.Hour),
	}

	res := uuid.NewString()

	tokensByAT[res] = newToken
	return newToken, res
}

func UserFromHttpRequest(r *http.Request) (*models.User, error) {

	raw := r.Header.Get("Authorization")
	if raw == "" {
		return nil, fmt.Errorf("Unauthorized")
	}
	if len(raw) < 7 || raw[:7] != "Bearer " {
		return nil, fmt.Errorf("Unauthorized")
	}
	token := raw[7:]

	if item, ok := tokensByAT[token]; ok {
		if item.accessTokenExpiresAt.After(time.Now()) {
			user, err := dbc.LoadUser(item.userRef)
			if err != nil {
				return nil, fmt.Errorf("Unauthorized")
			}
			if !user.IsActive() {
				return nil, fmt.Errorf("user isn't active")
			}
			return user, nil
		}
	}
	return nil, fmt.Errorf("Unauthorized")

}

type UserProfileResponse struct {
	Username   string `json:"userName"`
	Email      string `json:"email"`
	ShopManage bool   `json:"shopManage"`
	Admin      bool   `json:"admin"`
	Id         string `json:"id"`
}

func UserPublicProfile(w http.ResponseWriter, r *http.Request) {

	user, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}

	res := UserProfileResponse{
		Username:   user.Username(),
		Email:      user.Email(),
		ShopManage: user.ShopManager(),
		Admin:      user.Admin(),
		Id:         user.RefString(),
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
}

type TokenItem struct {
	userRef               string
	refreshToken          string
	accessTokenExpiresAt  time.Time
	refreshTokenExpiresAt time.Time
}

var tokensByAT = make(map[string]TokenItem)
