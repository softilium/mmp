package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
)

func initRouterAuth(router *http.ServeMux) {
	router.HandleFunc("/identity/register", UserRegister)
	router.HandleFunc("/identity/login", UserLogin)
	router.HandleFunc("/identity/logout", UserLogout)
	router.HandleFunc("/identity/myprofile", UserPublicProfile)
	router.HandleFunc("/identity/refresh", UserTokenRefresh)
}

type userPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileResponse struct {
	Username         string `json:"userName"`
	Email            string `json:"email"`
	ShopManage       bool   `json:"shopManage"`
	Admin            bool   `json:"admin"`
	Id               string `json:"id"`
	Description      string `json:"description"`
	TelegramUsername string `json:"telegramUsername"`
	BotChatId        int64  `json:"botChatId"`
}

type tokensResponse struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExporesAt  int64  `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt int64  `json:"refreshTokenExpiresAt"`
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

	newToken := models.GenerateToken(newUser)
	w.WriteHeader(http.StatusCreated)
	res := tokensResponse{
		AccessToken:           newToken.AccessToken,
		RefreshToken:          newToken.RefreshToken,
		AccessTokenExporesAt:  newToken.AccessTokenExpiresAt.UnixMilli(),
		RefreshTokenExpiresAt: newToken.RefreshTokenExpiresAt.UnixMilli(),
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
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

	newToken := models.GenerateToken(users[0])
	w.WriteHeader(http.StatusOK)
	res := tokensResponse{
		AccessToken:           newToken.AccessToken,
		RefreshToken:          newToken.RefreshToken,
		AccessTokenExporesAt:  newToken.AccessTokenExpiresAt.UnixMilli(),
		RefreshTokenExpiresAt: newToken.RefreshTokenExpiresAt.UnixMilli(),
	}
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

	user, err := models.UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}

	delete(models.TokensByAT, user.RefString())

	w.WriteHeader(http.StatusOK)
}

func UserPublicProfile(w http.ResponseWriter, r *http.Request) {

	user, err := models.UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}
	if r.Method == http.MethodGet {
		res := UserProfileResponse{
			Username:         user.Username(),
			Email:            user.Email(),
			ShopManage:       user.ShopManager(),
			Admin:            user.Admin(),
			Id:               user.RefString(),
			Description:      user.Description(),
			TelegramUsername: user.TelegramUsername(),
			//BotChatId:        user.BotChatId(),
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
	}
	if r.Method == http.MethodPut {
		payload := UserProfileResponse{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		user.SetUsername(payload.Username)
		user.SetEmail(payload.Email)
		user.SetShopManager(payload.ShopManage)
		user.SetAdmin(payload.Admin)
		user.SetDescription(payload.Description)
		user.SetTelegramUsername(payload.TelegramUsername)
		//user.SetBotChatId(payload.BotChatId)

		err = user.Save(r.Context())
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type refreshTokenPayload struct {
	RefreshToken string `json:"refreshToken"`
}

func UserTokenRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	payload := refreshTokenPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if payload.RefreshToken == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("refresh token is required"))
		return
	}
	for userKey, item := range models.TokensByAT {
		if item.RefreshToken == payload.RefreshToken && item.RefreshTokenExpiresAt.After(time.Now()) {
			user, err := models.Dbc.LoadUser(userKey)
			if err != nil {
				HandleErr(w, http.StatusUnauthorized, fmt.Errorf("user not found"))
				return
			}
			newToken := models.GenerateToken(user)
			models.TokensByAT[userKey] = newToken
			res := tokensResponse{
				AccessToken:           newToken.AccessToken,
				RefreshToken:          newToken.RefreshToken,
				AccessTokenExporesAt:  newToken.AccessTokenExpiresAt.UnixMilli(),
				RefreshTokenExpiresAt: newToken.RefreshTokenExpiresAt.UnixMilli(),
			}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				HandleErr(w, 0, err)
				return
			}
			return
		}
	}
	HandleErr(w, http.StatusUnauthorized, fmt.Errorf("invalid or expired refresh token"))
}
