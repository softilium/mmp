package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func HandleErr(w http.ResponseWriter, status int, err error) {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	errStr := ""
	if err != nil {
		log.Println(err.Error())
		errStr = err.Error()
	}
	http.Error(w, errStr, status)
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
	exist, _, err := models.Dbc.UserDef.SelectEntities(
		[]*elorm.Filter{elorm.AddFilterEQ(models.Dbc.UserDef.Email, payload.Email)}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(exist) > 0 {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("user with this email already exists: %s", payload.Email))
		return
	}
	newUser, err := models.Dbc.CreateUser()
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	newUser.SetUsername(payload.Email)
	newUser.SetEmail(payload.Email)
	newUser.SetPassword(payload.Password)
	err = newUser.Save(r.Context())
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
	users, _, err := models.Dbc.UserDef.SelectEntities(
		[]*elorm.Filter{elorm.AddFilterEQ(models.Dbc.UserDef.Email, payload.Email)}, nil, 1, 1)
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
	res := tokensResponse{AccessToken: at, RefreshToken: newToken.RefreshToken}
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

	user, ok := r.Context().Value(models.UserContextKey).(*models.User)
	if !ok {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}
	for token, item := range models.TokensByAT {
		if item.UserRef == user.RefString() {
			delete(models.TokensByAT, token)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func generateToken(user *models.User) (models.TokenItem, string) {
	newToken := models.TokenItem{
		UserRef:               user.RefString(),
		RefreshToken:          uuid.NewString(),
		AccessTokenExpiresAt:  time.Now().Add(5 * time.Minute),
		RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour),
	}

	res := uuid.NewString()

	models.TokensByAT[res] = newToken
	return newToken, res
}

func UserPublicProfile(w http.ResponseWriter, r *http.Request) {

	user, err := models.UserFromHttpRequest(r)
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

type UserProfileResponse struct {
	Username   string `json:"userName"`
	Email      string `json:"email"`
	ShopManage bool   `json:"shopManage"`
	Admin      bool   `json:"admin"`
	Id         string `json:"id"`
}
