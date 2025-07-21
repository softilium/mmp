package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/softilium/mmp-go/models"
)

func Migrate(w http.ResponseWriter, r *http.Request) {

	ctxMig := context.WithValue(context.Background(), "migration", true)

	user, err := models.UserFromHttpRequest(r)
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

	users, _, err := models.Dbc.UserDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to fetch users for migration: %v", err))
		return
	}

	for _, user := range users {
		if user.Username() == Cfg.AdminEmail {
			continue
		}
		err := models.Dbc.DeleteEntity(ctxMig, user.RefString())
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to delete user %s: %v", user.RefString(), err))
			return
		}
	}

	rows, err := mdb.Query(`
	
select 
"UserName", "Email", "Admin", "ShopManage", "TelegramCheckCode", "TelegramUserName", "TelegramVerified", "Description", "Id"
from "AspNetUsers" 
order by "Id"
	
	`)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to query users from migration database: %v", err))
		return
	}
	defer rows.Close()

	usersMap := make(map[int64]*models.User)

	for rows.Next() {

		var userName, email, telegramUserName, description string
		var admin, shopManage, telegramVerified bool
		var telegramCheckCode sql.NullString
		var id int64

		err := rows.Scan(&userName, &email, &admin, &shopManage, &telegramCheckCode, &telegramUserName, &telegramVerified, &description, &id)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan user row: %v", err))
			return
		}

		newUser, err := models.Dbc.CreateUser()
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
		err = newUser.Save(ctxMig)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to save user %s: %v", userName, err))
			return
		}

		usersMap[id] = newUser
	}

	// shops
	////////

	shops, _, err := models.Dbc.ShopDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to fetch shops for migration: %v", err))
		return
	}
	for _, shop := range shops {
		err := models.Dbc.DeleteEntity(ctxMig, shop.RefString())
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to delete shop %s: %v", shop.RefString(), err))
			return
		}
	}

	rows, err = mdb.Query(`
	
select 
"ID", "Caption", "Description", "DeliveryConditions", "CreatedByID", "CreatedOn", "ModifiedByID", "ModifiedOn", "DeletedByID", "DeletedOn" 
from "Shops"
	
	`)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to query shops from migration database: %v", err))
		return
	}
	defer rows.Close()

	for rows.Next() {

		ns, err := models.Dbc.CreateShop()
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to create new shop: %v", err))
			return
		}

		var id int64
		var createdById, modifiedById, deletedById sql.NullInt64

		err = rows.Scan(
			&id,
			ns.Values["Caption"],
			ns.Values["Description"],
			ns.Values["DeliveryConditions"],
			&createdById,
			ns.Values["CreatedAt"],
			&modifiedById,
			ns.Values["ModifiedAt"],
			&deletedById,
			ns.Values["DeletedAt"])
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan shop row: %v", err))
			return
		}

		if createdById.Valid {
			ns.SetCreatedBy(usersMap[createdById.Int64])
		}

		if modifiedById.Valid {
			ns.SetModifiedBy(usersMap[modifiedById.Int64])
		}
		if deletedById.Valid {
			ns.SetDeletedBy(usersMap[deletedById.Int64])
		}

		err = ns.Save(ctxMig)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to save shop %d: %v", id, err))
			return
		}

	}

}
