package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func Migrate(w http.ResponseWriter, r *http.Request) {

	user, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("unauthorized: %v", err))
		return
	}
	if !user.Admin() {
		HandleErr(w, http.StatusForbidden, fmt.Errorf("user isn't admin"))
		return
	}

	mdb, err := sql.Open("postgres", Cfg.DbConnectionStringForMigration)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to connect to the database for migration"))
		return
	}
	defer mdb.Close()

	// Users
	////////

	users, _, err := dbc.UserDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to fetch users for migration: %v", err))
		return
	}

	for _, user := range users {
		err := dbc.DeleteEntity(user.RefString())
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to delete user %s: %v", user.RefString(), err))
			return
		}
	}

	rows, err := mdb.Query(`
	
select 
"UserName", "Email", "Admin", "ShopManage", "TelegramCheckCode", "TelegramUserName", "TelegramVerified", "Description"
from "AspNetUsers" 
order by "Id"
	
	`)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to query users from migration database: %v", err))
		return
	}
	defer rows.Close()
	for rows.Next() {

		var userName, email, telegramUserName, description string
		var admin, shopManage, telegramVerified bool
		var telegramCheckCode sql.NullString

		err := rows.Scan(&userName, &email, &admin, &shopManage, &telegramCheckCode, &telegramUserName, &telegramVerified, &description)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan user row: %v", err))
			return
		}

		newUser, err := dbc.CreateUser()
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to create new user: %v", err))
			return
		}

		newUser.SetUsername(userName)
		newUser.SetPassword(Cfg.AdminPassword)
		newUser.SetEmail(email)
		newUser.SetAdmin(admin)
		newUser.SetShopManager(shopManage)
		newUser.SetTelegramCheckCode(telegramCheckCode.String)
		newUser.SetTelegramUsername(telegramUserName)
		newUser.SetTelegramVerified(telegramVerified)
		newUser.SetDescription(description)
		err = newUser.Save()
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to save user %s: %v", userName, err))
			return
		}
	}
}
