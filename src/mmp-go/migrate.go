package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
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

	usersMap := make(map[int64]*models.User)
	shopsMap := make(map[int64]*models.Shop)
	goodMap := make(map[int64]*models.Good)
	ordersMap := make(map[int64]*models.CustomerOrder)

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
u."UserName", u."Email", u."Admin", u."ShopManage", u."TelegramCheckCode", u."TelegramUserName", u."TelegramVerified", u."Description", u."Id", c."ChatId"
from "AspNetUsers" as u
left join "BotChats" as c on u."UserName"=c."UserName"
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
		var id int64
		var chatId sql.NullInt64

		err := rows.Scan(&userName, &email, &admin, &shopManage, &telegramCheckCode, &telegramUserName, &telegramVerified, &description, &id, &chatId)
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

		if chatId.Valid {
			newUser.SetBotChatId(chatId.Int64)
		}

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
from "Shops" order by "ID"
	
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

		ns.SetIsDeleted(deletedById.Valid && deletedById.Int64 > 0)

		err = ns.Save(ctxMig)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to save shop %d: %v", id, err))
			return
		}

		shopsMap[id] = ns

	}

	goods, _, err := models.Dbc.GoodDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to fetch goods for migration: %v", err))
		return
	}
	for _, good := range goods {
		err := models.Dbc.DeleteEntity(ctxMig, good.RefString())
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to delete good %s: %v", good.RefString(), err))
			return
		}
	}
	rows, err = mdb.Query(`

select 
"ID", "OwnerShopID", "Caption", "Article", "Url", "Description", "Price", "CreatedByID", 
"CreatedOn", "ModifiedByID", "ModifiedOn", "DeletedByID", "OrderInShop"
from "Goods" order by "ID"

`)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to query goods from migration database: %v", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		ng, err := models.Dbc.CreateGood()
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to create new good: %v", err))
			return
		}

		var id, ownerShopId int64
		var createdById, modifiedById, deletedById sql.NullInt64

		err = rows.Scan(
			&id,
			&ownerShopId,
			ng.Values["Caption"],
			ng.Values["Article"],
			ng.Values["Url"],
			ng.Values["Description"],
			ng.Values["Price"],
			&createdById,
			ng.Values["CreatedAt"],
			&modifiedById,
			ng.Values["ModifiedAt"],
			&deletedById,
			ng.Values["OrderInShop"])
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan good row: %v", err))
			return
		}

		if createdById.Valid {
			ng.SetCreatedBy(usersMap[createdById.Int64])
		}

		if modifiedById.Valid {
			ng.SetModifiedBy(usersMap[modifiedById.Int64])
		}

		if deletedById.Valid {
			ng.SetDeletedBy(usersMap[deletedById.Int64])
		}

		if ownerShopId > 0 {
			ng.SetOwnerShop(shopsMap[ownerShopId])
		}

		ng.SetIsDeleted(deletedById.Valid && deletedById.Int64 > 0)

		err = ng.Save(ctxMig)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan good row: %v", err))
			return
		}

		goodMap[id] = ng

	}

	// images
	/////////

	items, _ := os.ReadDir(Cfg.ImagesFolderForMigration)
	for _, item := range items {

		fn := item.Name()
		tokens := strings.Split(fn, "-")
		if len(tokens) != 3 {
			HandleErr(w, 0, fmt.Errorf("invalid image file name %s", fn))
			return
		}
		if tokens[0] != "goodImage" {
			HandleErr(w, 0, fmt.Errorf("invalid image file name %s, expected 'goodImage-<goodId>-<imageNum>'", fn))
			return
		}

		t1, err := strconv.Atoi(tokens[1])
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("invalid good id in image file name %s: %v", fn, err))
			return
		}

		g, ok := goodMap[int64(t1)]
		if !ok {
			fmt.Printf("good with id %s not found\n", tokens[1])
			continue
		}
		tokens[1] = g.RefString()
		newFn := strings.Join(tokens, "-")

		data, err := os.ReadFile(Cfg.ImagesFolderForMigration + "/" + fn)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to read image file %s: %v", fn, err))
			return
		}

		err = os.WriteFile(Cfg.ImagesFolder+"/"+newFn, data, 0644)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to read image file %s: %v", fn, err))
			return
		}
	}

	// orders
	/////////

	orders, _, err := models.Dbc.CustomerOrderDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to fetch orders for migration: %v", err))
		return
	}
	for _, order := range orders {
		err := models.Dbc.DeleteEntity(ctxMig, order.RefString())
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to delete order %s: %v", order.RefString(), err))
			return
		}
	}
	rows, err = mdb.Query(`

select 
	"ID", "SenderID", "Status", "Qty", "Sum", "CreatedByID", "CreatedOn", "ModifiedByID", "ModifiedOn", 
	"IsDeleted", "DeletedByID", "DeletedOn", "CustomerComment", "ExpectedDeliveryDate", "SenderComment" 
from "Orders"
order by "ID"

	`)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to query orders from migration database: %v", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		no, err := models.Dbc.CreateCustomerOrder()
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to create new order: %v", err))
			return
		}

		var id, senderId int64
		var createdById, modifiedById, deletedById sql.NullInt64
		var status int64

		err = rows.Scan(
			&id,
			&senderId,
			&status,
			no.Values["Qty"],
			no.Values["Sum"],
			&createdById,
			no.Values["CreatedAt"],
			&modifiedById,
			no.Values["ModifiedAt"],
			no.Values["IsDeleted"],
			&deletedById,
			no.Values["DeletedAt"],
			no.Values["CustomerComment"],
			no.Values["ExpectedDeliveryDate"],
			no.Values["SenderComment"])
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan order row: %v", err))
			return
		}

		if createdById.Valid {
			no.SetCreatedBy(usersMap[createdById.Int64])
		}

		if modifiedById.Valid {
			no.SetModifiedBy(usersMap[modifiedById.Int64])
		}

		if deletedById.Valid {
			no.SetDeletedBy(usersMap[deletedById.Int64])
		}

		if senderId > 0 {
			no.SetSender(usersMap[senderId])
		}

		no.SetStatus(status)

		no.SetIsDeleted(deletedById.Valid && deletedById.Int64 > 0)

		err = no.Save(ctxMig)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to save order %d: %v", id, err))
			return
		}
		ordersMap[id] = no
	}

	// order lines
	///////////

	orderLines, _, err := models.Dbc.OrderLineDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to fetch order lines for migration: %v", err))
		return
	}
	for _, orderLine := range orderLines {
		err := models.Dbc.DeleteEntity(ctxMig, orderLine.RefString())
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to delete order line %s: %v", orderLine.RefString(), err))
			return
		}
	}

	rows, err = mdb.Query(`

select 
"ID", "ShopID", "OrderID", "GoodID", "Qty", "Sum", "CreatedByID", "CreatedOn", "ModifiedByID", "ModifiedOn", "IsDeleted", "DeletedByID", "DeletedOn"
from "OrderLines"
order by "ID"

	`)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("failed to query order lines from migration database: %v", err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		nol, err := models.Dbc.CreateOrderLine()
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to create new order line: %v", err))
			return
		}

		var id, shopId, orderId, goodId int64
		var createdById, modifiedById, deletedById sql.NullInt64

		err = rows.Scan(
			&id,
			&shopId,
			&orderId,
			&goodId,
			nol.Values["Qty"],
			nol.Values["Sum"],
			&createdById,
			nol.Values["CreatedAt"],
			&modifiedById,
			nol.Values["ModifiedAt"],
			nol.Values["IsDeleted"],
			&deletedById,
			nol.Values["DeletedAt"])
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to scan order line row: %v", err))
			return
		}

		if createdById.Valid {
			nol.SetCreatedBy(usersMap[createdById.Int64])
		}

		if modifiedById.Valid {
			nol.SetModifiedBy(usersMap[modifiedById.Int64])
		}

		if deletedById.Valid {
			nol.SetDeletedBy(usersMap[deletedById.Int64])
		}

		if shopId > 0 {
			nol.SetShop(shopsMap[shopId])
		}

		if orderId > 0 {
			nol.SetCustomerOrder(ordersMap[orderId])
		}

		if goodId > 0 {
			nol.SetGood(goodMap[goodId])
		}

		nol.SetIsDeleted(deletedById.Valid && deletedById.Int64 > 0)

		err = nol.Save(ctxMig)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("failed to save order line %d: %v", id, err))
			return
		}
	}

	models.TokensByAT = make(map[string]models.TokenItem)
	imagesCache = expirable.NewLRU[string, *[]byte](1000, nil, time.Minute*60*24)
	thumbsCache = expirable.NewLRU[string, *[]byte](1000, nil, time.Minute*60*24)

}
