package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
)

func initRouterOrders(router *http.ServeMux) {

	router.HandleFunc("/api/orders/checkout", checkoutHandler)

	router.HandleFunc("/api/orders/statuses", func(w http.ResponseWriter, r *http.Request) {
		res := map[int]string{
			models.OrderStatusNew:            "New",
			models.OrderStatusInProcess:      "In Process",
			models.OrderStatusReadyToDeliver: "Ready to Deliver",
			models.OrderStatusDelivering:     "Delivering",
			models.OrderStatusDone:           "Done",
			models.OrderStatusCanceled:       "Canceled",
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
	})
	outbox := elorm.CreateStdRestApiConfig(
		*models.Dbc.CustomerOrderDef.EntityDef,
		models.Dbc.LoadCustomerOrder,
		models.Dbc.CustomerOrderDef.SelectEntities,
		models.Dbc.CreateCustomerOrder)
	outbox.EnablePost = false
	outbox.EnablePut = false
	outbox.EnableDelete = false
	outbox.BeforeMiddleware = UserRequired
	outbox.Context = models.HttpUserContext
	outbox.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		user, err := models.UserFromHttpRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error getting user from request: %w", err)
		}
		showAll := r.URL.Query().Get("showAll")
		res := []*elorm.Filter{}
		if showAll != "1" {
			res = append(res, elorm.AddFilterNOTIN(models.Dbc.CustomerOrderDef.Status, models.OrderStatusCanceled, models.OrderStatusDone))
		}
		res = append(res, elorm.AddFilterEQ(models.Dbc.CustomerOrderDef.IsDeleted, false))
		res = append(res, elorm.AddFilterEQ(models.Dbc.CustomerOrderDef.CreatedBy, user.RefString()))
		return res, nil
	}
	outbox.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: models.Dbc.CustomerOrderDef.CreatedAt, Asc: false}}, nil
	}
	outbox.Context = models.HttpUserContext
	router.HandleFunc("/api/orders/outbox", elorm.HandleRestApi(outbox))

	inbox := elorm.CreateStdRestApiConfig(
		*models.Dbc.CustomerOrderDef.EntityDef,
		models.Dbc.LoadCustomerOrder,
		models.Dbc.CustomerOrderDef.SelectEntities,
		models.Dbc.CreateCustomerOrder)
	inbox.EnablePost = false
	inbox.EnablePut = false
	inbox.EnableDelete = false
	inbox.BeforeMiddleware = UserRequired
	inbox.Context = models.HttpUserContext
	inbox.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		user, err := models.UserFromHttpRequest(r)
		if err != nil {
			return nil, fmt.Errorf("error getting user from request: %w", err)
		}
		showAll := r.URL.Query().Get("showAll")
		res := []*elorm.Filter{}
		if showAll != "1" {
			res = append(res, elorm.AddFilterNOTIN(models.Dbc.CustomerOrderDef.Status, models.OrderStatusCanceled, models.OrderStatusDone))
		}
		res = append(res, elorm.AddFilterEQ(models.Dbc.CustomerOrderDef.IsDeleted, false))
		res = append(res, elorm.AddFilterEQ(models.Dbc.CustomerOrderDef.Sender, user.RefString()))
		return res, nil
	}
	inbox.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: models.Dbc.CustomerOrderDef.CreatedAt, Asc: false}}, nil
	}
	inbox.Context = models.HttpUserContext
	router.HandleFunc("/api/orders/inbox", elorm.HandleRestApi(inbox))

	ordersGetterEditor := elorm.CreateStdRestApiConfig(
		*models.Dbc.CustomerOrderDef.EntityDef,
		models.Dbc.LoadCustomerOrder,
		models.Dbc.CustomerOrderDef.SelectEntities,
		models.Dbc.CreateCustomerOrder)
	ordersGetterEditor.EnablePost = false
	ordersGetterEditor.EnablePut = true
	ordersGetterEditor.EnableDelete = false
	ordersGetterEditor.EnableGetList = true
	ordersGetterEditor.EnableGetOne = true
	ordersGetterEditor.BeforeMiddleware = UserRequired
	ordersGetterEditor.Context = models.HttpUserContext
	inbox.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: models.Dbc.CustomerOrderDef.CreatedAt, Asc: false}}, nil
	}
	ordersGetterEditor.Context = models.HttpUserContext
	router.HandleFunc("/api/orders", elorm.HandleRestApi(ordersGetterEditor))

	lines := elorm.CreateStdRestApiConfig(
		*models.Dbc.OrderLineDef.EntityDef,
		models.Dbc.LoadOrderLine,
		models.Dbc.OrderLineDef.SelectEntities,
		models.Dbc.CreateOrderLine)
	lines.EnablePost = false
	lines.EnablePut = false
	lines.EnableDelete = false
	lines.EnableGetList = true
	lines.EnableGetOne = false
	lines.BeforeMiddleware = UserRequired
	lines.Context = models.HttpUserContext
	lines.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		orderref := r.URL.Query().Get("orderref")
		if orderref == "" {
			return nil, fmt.Errorf("orderref is required")
		}
		res := []*elorm.Filter{
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CustomerOrder, orderref),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.IsDeleted, false),
		}
		return res, nil
	}
	lines.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: models.Dbc.OrderLineDef.CreatedAt, Asc: false}}, nil
	}
	router.HandleFunc("/api/orderlines", elorm.HandleRestApi(lines))

}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	if !UserRequired(w, r) {
		return
	}
	user, err := models.UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error getting user from request: %w", err))
		return
	}

	sender := r.URL.Query().Get("sender")
	if sender == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("sender is required"))
		return
	}
	senderObj, err := models.Dbc.LoadUser(sender)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error loading sender user: %w", err))
		return
	}

	customerComment := ""
	body, err := io.ReadAll(r.Body)
	if err == nil {
		customerComment = string(body)
	}

	newLines, _, err := models.Dbc.OrderLineDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.IsDeleted, false),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CustomerOrder, ""),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CreatedBy, user.RefString()),
		}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error selecting order lines: %w", err))
		return
	}

	newOrder, err := models.Dbc.CreateCustomerOrder()
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error creating new order: %w", err))
		return
	}
	newOrder.SetSender(senderObj)
	newOrder.SetCustomerComment(customerComment)

	ctx := models.AddUserContext(r.Context(), user)

	tx, err := models.Dbc.BeginTran()
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error starting transaction: %w", err))
		return
	}
	defer models.Dbc.RollbackTran(tx)

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
	newOrder.SetStatus(models.OrderStatusNew)
	err = newOrder.Save(ctx)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error saving new order: %w", err))
		return
	}
	err = models.Dbc.CommitTran(tx)
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
