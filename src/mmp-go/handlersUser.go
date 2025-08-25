package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/softilium/elorm"
)

func initRouterAuth(router *http.ServeMux) {
	router.HandleFunc("/identity/register", UserRegister)
	router.HandleFunc("/identity/login", UserLogin)
	router.HandleFunc("/identity/logout", UserLogout)
	router.HandleFunc("/identity/myprofile", UserMyProfile)
	router.HandleFunc("/identity/profiles", UserPublicProfile)
	router.HandleFunc("/identity/refresh", RefreshToken)
	router.HandleFunc("/api/profiles/sendmsg", SendMessageToUser)
	router.HandleFunc("/api/users/sendreset", SendReset)
	router.HandleFunc("/api/users/resetpwd", ResetPwd)
}

type userPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileResponse struct {
	Username          string `json:"userName"`
	Email             string `json:"email"`
	ShopManage        bool   `json:"shopManage"`
	Admin             bool   `json:"admin"`
	Id                string `json:"id"`
	Description       string `json:"description"`
	TelegramUsername  string `json:"telegramUsername"`
	BotChatId         int64  `json:"botChatId"`
	TelegramVerified  bool   `json:"telegramVerified"`
	TelegramCheckCode string `json:"telegramCheckCode"`
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
	exist, _, err := DB.UserDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(DB.UserDef.Email, payload.Email),
		}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(exist) > 0 {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("user with this email already exists: %s", payload.Email))
		return
	}
	newUser, err := DB.CreateUser()
	if err != nil {
		HandleErr(w, 0, err)
		return
	}

	passwordHash, err := HashPassword(payload.Password)
	if err != nil {
		HandleErr(w, 0, err)
	}

	newUser.SetUsername(payload.Email)
	newUser.SetEmail(payload.Email)
	newUser.SetPasswordHash(passwordHash)
	err = newUser.Save(r.Context())
	if err != nil {
		HandleErr(w, 0, err)
		return
	}

	newToken, err := GenerateToken(newUser)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	res := tokensResponse{
		AccessToken:           newToken.AccessToken(),
		RefreshToken:          newToken.RefreshToken(),
		AccessTokenExporesAt:  newToken.AccessTokenExpiresAt().UnixMilli(),
		RefreshTokenExpiresAt: newToken.RefreshTokenExpiresAt().UnixMilli(),
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
		HandleErr(w, http.StatusBadRequest, err)
		return
	}
	users, _, err := DB.UserDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(DB.UserDef.Email, payload.Email),
		}, nil, 1, 1)
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

	if !CheckPasswordHash(payload.Password, users[0].PasswordHash()) {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	newToken, err := GenerateToken(users[0])
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	res := tokensResponse{
		AccessToken:           newToken.AccessToken(),
		RefreshToken:          newToken.RefreshToken(),
		AccessTokenExporesAt:  newToken.AccessTokenExpiresAt().UnixMilli(),
		RefreshTokenExpiresAt: newToken.RefreshTokenExpiresAt().UnixMilli(),
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

	_, token, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return
	}
	err = DB.DeleteEntity(r.Context(), token.RefString())
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UserMyProfile(w http.ResponseWriter, r *http.Request) {

	user, _, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return
	}
	if r.Method == http.MethodGet {
		res := UserProfileResponse{
			Username:          user.Username(),
			Email:             user.Email(),
			ShopManage:        user.ShopManager(),
			Admin:             user.Admin(),
			Id:                user.RefString(),
			Description:       user.Description(),
			TelegramUsername:  user.TelegramUsername(),
			BotChatId:         user.TelegramChatId(),
			TelegramVerified:  user.TelegramVerified(),
			TelegramCheckCode: user.TelegramCheckCode(),
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
			HandleErr(w, http.StatusBadRequest, err)
			return
		}
		user.SetUsername(payload.Username)
		user.SetEmail(payload.Email)
		user.SetShopManager(payload.ShopManage)
		user.SetAdmin(payload.Admin)
		user.SetDescription(payload.Description)
		user.SetTelegramUsername(payload.TelegramUsername)
		//user.SetTelegramVerified(payload.TelegramVerified)
		//user.SetTelegramCheckCode(payload.TelegramCheckCode)

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

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	payload := refreshTokenPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		HandleErr(w, http.StatusBadRequest, err)
		return
	}
	if payload.RefreshToken == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("refresh token is required"))
		return
	}
	dbrt, _, err := DB.TokenDef.SelectEntities([]*elorm.Filter{
		elorm.AddFilterEQ(DB.TokenDef.RefreshToken, payload.RefreshToken),
		elorm.AddFilterGT(DB.TokenDef.RefreshTokenExpiresAt, time.Now()),
	}, nil, 1, 1)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(dbrt) == 1 {
		user := dbrt[0].User()
		if !user.IsActive() {
			HandleErr(w, http.StatusUnauthorized, fmt.Errorf("user isn't active"))
			return
		}
		newToken, err := GenerateToken(user)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		res := tokensResponse{
			AccessToken:           newToken.AccessToken(),
			RefreshToken:          newToken.RefreshToken(),
			AccessTokenExporesAt:  newToken.AccessTokenExpiresAt().UnixMilli(),
			RefreshTokenExpiresAt: newToken.RefreshTokenExpiresAt().UnixMilli(),
		}
		err = DB.DeleteEntity(context.Background(), dbrt[0].RefString())
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		return

	}
	HandleErr(w, http.StatusUnauthorized, fmt.Errorf("refresh token is expired or invalid"))
}

func UserPublicProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}

	userref := r.URL.Query().Get("userref")
	if userref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("userref is required"))
		return
	}

	user, err := DB.LoadUser(userref)
	if err != nil {
		HandleErr(w, http.StatusNotFound, fmt.Errorf("user not found: %s", userref))
		return
	}

	res := UserProfileResponse{
		Username:         user.Username(),
		Email:            user.Email(),
		ShopManage:       user.ShopManager(),
		Admin:            user.Admin(),
		Id:               user.RefString(),
		Description:      user.Description(),
		TelegramUsername: user.TelegramUsername(),
		//BotChatId:        user.BotChatId(),
		TelegramVerified: user.TelegramVerified(),
		//TelegramCheckCode: user.TelegramCheckCode(),
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
}

func SendMessageToUser(w http.ResponseWriter, r *http.Request) {
	msg := ""
	body, err := io.ReadAll(r.Body)
	if err == nil {
		msg = string(body)
	}
	if msg == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("message is required"))
		return
	}
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	sender, _, _ := UserFromHttpRequest(r)
	if sender != nil {
		msg = fmt.Sprintf("Сообщение от %s:\n\r\n\r%s", sender.Username(), msg)
	}
	uref := r.URL.Query().Get("userref")
	if uref == "" {
		admins, _, err := DB.UserDef.SelectEntities([]*elorm.Filter{elorm.AddFilterEQ(DB.UserDef.Admin, true)}, nil, 0, 0)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		for _, admin := range admins {
			if admin.TelegramVerified() && Bot != nil {
				_, _ = Bot.Send(tgbotapi.NewMessage(admin.TelegramChatId(), msg))
			}
		}
	} else {
		user, err := DB.LoadUser(uref)
		if err != nil {
			HandleErr(w, http.StatusNotFound, fmt.Errorf("user not found: %s", uref))
			return
		}
		if user.TelegramVerified() && Bot != nil {
			_, _ = Bot.Send(tgbotapi.NewMessage(user.TelegramChatId(), msg))
		} else {
			HandleErr(w, http.StatusBadRequest, fmt.Errorf("user %s is not verified in Telegram", uref))
			return
		}
	}
}

func SendReset(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("email is required"))
		return
	}

	users, _, err := DB.UserDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(DB.UserDef.Email, email),
		}, nil, 1, 1)

	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(users) != 1 {
		HandleErr(w, http.StatusNotFound, fmt.Errorf("user with this email not found"))
		return
	}
	user := users[0]
	user.SetRecoverCodeEmail(elorm.NewRef())
	user.SetRecoverCodeDeadline(time.Now().Add(1 * time.Hour))
	err = user.Save(r.Context())
	if err != nil {
		HandleErr(w, 0, err)
		return
	}

	err = SendEmail(
		user.Email(),
		"Сброс пароля для river-stores.com",
		fmt.Sprintf("Вы запросили сброс пароля на сайте river-stores.com\n\rПерейдите по ссылке https://river-stores.com/resetpwd?code=%s в течение часа с момента отправки запроса на сброс пароля.", user.RecoverCodeEmail()))

	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func ResetPwd(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}

	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")
	pwd := r.URL.Query().Get("pwd")
	if email == "" || code == "" || pwd == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("email, code and pwd are required"))
		return
	}

	users, _, err := DB.UserDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(DB.UserDef.Email, email),
			elorm.AddFilterEQ(DB.UserDef.RecoverCodeEmail, code),
			elorm.AddFilterGT(DB.UserDef.RecoverCodeDeadline, time.Now()),
		}, nil, 1, 1)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(users) != 1 {
		HandleErr(w, http.StatusNotFound, fmt.Errorf("user with this email and code not found or code expired"))
		return
	}
	user := users[0]
	passwordHash, err := HashPassword(pwd)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	user.SetPasswordHash(passwordHash)
	user.SetRecoverCodeEmail("")
	err = user.Save(r.Context())
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
