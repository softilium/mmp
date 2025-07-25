package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
)

func initRouterOrders(router *http.ServeMux) {

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

}
