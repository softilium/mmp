package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
)

func initRouterBasket(router *http.ServeMux) {

	basketRestApiConfig := elorm.CreateStdRestApiConfig(
		*models.Dbc.OrderLineDef.EntityDef,
		models.Dbc.LoadOrderLine,
		models.Dbc.OrderLineDef.SelectEntities,
		models.Dbc.CreateOrderLine)
	basketRestApiConfig.EnablePost = false
	basketRestApiConfig.BeforeMiddleware = func(w http.ResponseWriter, r *http.Request) bool {
		_, err := models.UserFromHttpRequest(r)
		if err != nil {
			HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
			return false
		}
		return true
	}
	basketRestApiConfig.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		user, err := models.UserFromHttpRequest(r)
		if err != nil {
			return nil, err
		}
		res := []*elorm.Filter{}
		res = append(res, elorm.AddFilterEQ(models.Dbc.OrderLineDef.IsDeleted, false))
		res = append(res, elorm.AddFilterEQ(models.Dbc.OrderLineDef.CustomerOrder, "")) // only active basket lines (order is empty)
		res = append(res, elorm.AddFilterEQ(models.Dbc.OrderLineDef.CreatedBy, user.RefString()))
		return res, nil
	}
	basketRestApiConfig.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: models.Dbc.OrderLineDef.Ref, Asc: true}}, nil
	}
	basketRestApiConfig.Context = models.HttpUserContext
	router.HandleFunc("/api/basket", elorm.HandleRestApi(basketRestApiConfig))

	router.HandleFunc("/api/basket/increase", increaseBasket)
	router.HandleFunc("/api/basket/decrease", decreaseBasket)
	router.HandleFunc("/api/basket/merge", mergeBasket)

}

func increaseBasket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	user, err := models.UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}
	gref := r.URL.Query().Get("goodref")
	if gref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("goodref is required"))
		return
	}
	good, err := models.Dbc.LoadGood(gref)
	if err != nil {
		HandleErr(w, http.StatusNotFound, fmt.Errorf("good not found: %v", err))
		return
	}
	exists, _, err := models.Dbc.OrderLineDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.IsDeleted, false),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CustomerOrder, ""), // only active basket lines (order is empty)
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CreatedBy, user.RefString()),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.Good, gref),
		}, nil, 0, 0)
	if err != nil {
		HandleErr(w, http.StatusNotFound, err)
		return
	}
	if len(exists) > 1 {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("more than one basket line found for good %s", gref))
		return
	}
	if len(exists) == 1 {
		exists[0].SetQty(exists[0].Qty() + 1)
		exists[0].SetSum(exists[0].Sum() + good.Price())
		err = exists[0].Save(r.Context())
		if err != nil {
			HandleErr(w, 0, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	newLine, err := models.Dbc.CreateOrderLine()
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	newLine.SetGood(good)
	newLine.SetQty(1)
	newLine.SetSum(good.Price())
	newLine.SetShop(good.OwnerShop())
	newLine.SetCreatedBy(user)
	err = newLine.Save(r.Context())
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func decreaseBasket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleErr(w, http.StatusMethodNotAllowed, nil)
		return
	}
	user, err := models.UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}
	gref := r.URL.Query().Get("goodref")
	if gref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("goodref is required"))
		return
	}
	good, err := models.Dbc.LoadGood(gref)
	if err != nil {
		HandleErr(w, http.StatusNotFound, fmt.Errorf("good not found: %v", err))
		return
	}
	exists, _, err := models.Dbc.OrderLineDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.IsDeleted, false),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CustomerOrder, ""), // only active basket lines (order is empty)
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CreatedBy, user.RefString()),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.Good, gref),
		}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, err)
		return
	}
	if len(exists) > 1 {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("more than one basket line found for good %s", gref))
		return
	}
	if len(exists) > 0 {
		exists[0].SetQty(exists[0].Qty() - 1)
		exists[0].SetSum(exists[0].Qty() * good.Price())
		if exists[0].Qty() < 1 {
			err = models.Dbc.DeleteEntity(r.Context(), exists[0].RefString())
			if err != nil {
				HandleErr(w, 0, err)
				return
			}
		} else {
			err = exists[0].Save(r.Context())
			if err != nil {
				HandleErr(w, 0, err)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

type MergeBasketItem struct {
	GoodId   string  `json:"goodId"`
	Quantity float64 `json:"quantity"`
}

func mergeBasket(w http.ResponseWriter, r *http.Request) {

	user, err := models.UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, nil)
		return
	}

	dbCtx := models.AddUserContext(r.Context(), user)

	var newItems []MergeBasketItem
	err = json.NewDecoder(r.Body).Decode(&newItems)
	if err != nil {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}

	existItems, _, err := models.Dbc.OrderLineDef.SelectEntities(
		[]*elorm.Filter{
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.IsDeleted, false),
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CustomerOrder, ""), // only active basket lines (order is empty)
			elorm.AddFilterEQ(models.Dbc.OrderLineDef.CreatedBy, user.RefString()),
		}, nil, 0, 0)
	if err != nil {
		HandleErr(w, 0, fmt.Errorf("error selecting basket lines: %w", err))
		return
	}

	for _, newItem := range newItems {
		updated := false
		for _, existItem := range existItems {
			if existItem.Good().RefString() == newItem.GoodId {
				// Update existing item
				existItem.SetQty(existItem.Qty() + newItem.Quantity)
				existItem.SetSum(existItem.Qty() * existItem.Good().Price())
				err = existItem.Save(dbCtx)
				if err != nil {
					HandleErr(w, 0, fmt.Errorf("error updating existing item: %w", err))
					return
				}
				updated = true
				break
			}
		}
		if !updated {
			// Create new item
			good, err := models.Dbc.LoadGood(newItem.GoodId)
			if err != nil {
				HandleErr(w, http.StatusNotFound, fmt.Errorf("good not found: %v", err))
				return
			}
			newLine, err := models.Dbc.CreateOrderLine()
			if err != nil {
				HandleErr(w, 0, fmt.Errorf("error creating new order line: %w", err))
				return
			}
			newLine.SetGood(good)
			newLine.SetQty(newItem.Quantity)
			newLine.SetSum(newItem.Quantity * good.Price())
			newLine.SetShop(good.OwnerShop())
			newLine.SetCreatedBy(user)
			err = newLine.Save(dbCtx)
			if err != nil {
				HandleErr(w, 0, fmt.Errorf("error saving new order line: %w", err))
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
