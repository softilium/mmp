package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/softilium/elorm"
)

func initRouterOrders(router *http.ServeMux) {

	router.HandleFunc("/api/orders/checkout", checkoutHandler)

	router.HandleFunc("/api/orders/statuses", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(OrderStatusCaptions)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
	})
	outbox := elorm.CreateStdRestApiConfig(
		*DB.CustomerOrderDef.EntityDef,
		DB.LoadCustomerOrder,
		DB.CustomerOrderDef.SelectEntities,
		DB.CreateCustomerOrder)
	outbox.DefaultPageSize = 0
	outbox.EnablePost = false
	outbox.EnablePut = false
	outbox.EnableDelete = false
	outbox.BeforeMiddleware = UserRequired
	outbox.Context = LoadUserFromHttpToContext
	outbox.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		user, _, err := UserFromHttpRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error getting user from request: %w", err)
		}
		showAll := r.URL.Query().Get("showAll")
		res := []*elorm.Filter{}
		if showAll != "1" {
			res = append(res, elorm.AddFilterNOTIN(DB.CustomerOrderDef.Status, OrderStatusCanceled, OrderStatusDone))
		}
		res = append(res, elorm.AddFilterEQ(DB.CustomerOrderDef.IsDeleted, false))
		res = append(res, elorm.AddFilterEQ(DB.CustomerOrderDef.CreatedBy, user.RefString()))
		return res, nil
	}
	outbox.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: DB.CustomerOrderDef.CreatedAt, Asc: false}}, nil
	}
	outbox.Context = LoadUserFromHttpToContext
	router.HandleFunc("/api/orders/outbox", elorm.HandleRestApi(outbox))

	inbox := elorm.CreateStdRestApiConfig(
		*DB.CustomerOrderDef.EntityDef,
		DB.LoadCustomerOrder,
		DB.CustomerOrderDef.SelectEntities,
		DB.CreateCustomerOrder)
	inbox.DefaultPageSize = 0
	inbox.EnablePost = false
	inbox.EnablePut = false
	inbox.EnableDelete = false
	inbox.BeforeMiddleware = UserRequired
	inbox.Context = LoadUserFromHttpToContext
	inbox.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		user, _, err := UserFromHttpRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error getting user from request: %w", err)
		}
		showAll := r.URL.Query().Get("showAll")
		res := []*elorm.Filter{}
		if showAll != "1" {
			res = append(res, elorm.AddFilterNOTIN(DB.CustomerOrderDef.Status, OrderStatusCanceled, OrderStatusDone))
		}
		res = append(res, elorm.AddFilterEQ(DB.CustomerOrderDef.IsDeleted, false))
		res = append(res, elorm.AddFilterEQ(DB.CustomerOrderDef.Sender, user.RefString()))
		return res, nil
	}
	inbox.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: DB.CustomerOrderDef.CreatedAt, Asc: false}}, nil
	}
	inbox.Context = LoadUserFromHttpToContext
	router.HandleFunc("/api/orders/inbox", elorm.HandleRestApi(inbox))

	ordersGetterEditor := elorm.CreateStdRestApiConfig(
		*DB.CustomerOrderDef.EntityDef,
		DB.LoadCustomerOrder,
		DB.CustomerOrderDef.SelectEntities,
		DB.CreateCustomerOrder)
	ordersGetterEditor.DefaultPageSize = 0
	ordersGetterEditor.EnablePost = false
	ordersGetterEditor.EnablePut = true
	ordersGetterEditor.EnableDelete = false
	ordersGetterEditor.EnableGetList = true
	ordersGetterEditor.EnableGetOne = true
	ordersGetterEditor.BeforeMiddleware = UserRequired
	ordersGetterEditor.Context = LoadUserFromHttpToContext
	inbox.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: DB.CustomerOrderDef.CreatedAt, Asc: false}}, nil
	}
	ordersGetterEditor.Context = LoadUserFromHttpToContext
	router.HandleFunc("/api/orders", elorm.HandleRestApi(ordersGetterEditor))

	lines := elorm.CreateStdRestApiConfig(
		*DB.OrderLineDef.EntityDef,
		DB.LoadOrderLine,
		DB.OrderLineDef.SelectEntities,
		DB.CreateOrderLine)
	lines.DefaultPageSize = 0
	lines.EnablePost = false
	lines.EnablePut = false
	lines.EnableDelete = false
	lines.EnableGetList = true
	lines.EnableGetOne = false
	lines.BeforeMiddleware = UserRequired
	lines.Context = LoadUserFromHttpToContext
	lines.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		orderref := r.URL.Query().Get("orderref")
		if orderref == "" {
			return nil, fmt.Errorf("orderref is required")
		}
		res := []*elorm.Filter{
			elorm.AddFilterEQ(DB.OrderLineDef.CustomerOrder, orderref),
			elorm.AddFilterEQ(DB.OrderLineDef.IsDeleted, false),
		}
		return res, nil
	}
	lines.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: DB.OrderLineDef.CreatedAt, Asc: false}}, nil
	}
	router.HandleFunc("/api/orderlines", elorm.HandleRestApi(lines))

}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	if !UserRequired(w, r) {
		return
	}
	user, _, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("error getting user from request: %w", err))
		return
	}

	sender := r.URL.Query().Get("sender")
	if sender == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("sender is required"))
		return
	}
	senderObj, err := DB.LoadUser(sender)
	if err != nil {
		HandleErr(w, http.StatusNotFound, fmt.Errorf("error loading sender user: %w", err))
		return
	}

	customerComment := ""
	body, err := io.ReadAll(r.Body)
	if err == nil {
		customerComment = string(body)
	}

	newLines, _, err := DB.OrderLineDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(DB.OrderLineDef.IsDeleted, false),
			elorm.AddFilterEQ(DB.OrderLineDef.CustomerOrder, ""),
			elorm.AddFilterEQ(DB.OrderLineDef.CreatedBy, user.RefString()),
		}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error selecting order lines: %w", err))
		return
	}

	newOrder, err := DB.CreateCustomerOrder()
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error creating new order: %w", err))
		return
	}
	newOrder.SetSender(senderObj)
	newOrder.SetCustomerComment(customerComment)

	ctx := AddUserContext(r.Context(), user)

	tx, err := DB.BeginTran()
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error starting transaction: %w", err))
		return
	}
	defer DB.RollbackTran(tx)

	for _, line := range newLines {
		if line.Good().CreatedBy().RefString() == senderObj.RefString() {
			line.SetCustomerOrder(newOrder)
			err := line.Save(ctx)
			if err != nil {
				HandleErr(w, 0, fmt.Errorf("error saving order line: %w", err))
				return
			}
			newOrder.SetQty(newOrder.Qty() + line.Qty())
			newOrder.SetSum(newOrder.Sum() + line.Sum())
		}
	}
	newOrder.SetStatus(OrderStatusNew)
	err = newOrder.Save(ctx)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error saving new order: %w", err))
		return
	}
	err = DB.CommitTran(tx)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error committing transaction: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newOrder)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error encoding new order: %w", err))
		return
	}
}
