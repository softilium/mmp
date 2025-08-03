package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/softilium/elorm"
)

func initRouterSearchGoods(router *http.ServeMux) {
	router.HandleFunc("/api/goods/search", GoodsSearch)
}

func GoodsSearch(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	q = strings.TrimSpace(q)
	q = strings.ReplaceAll(q, "  ", " ")

	tokens := strings.Split(q, " ")

	gd := DB.GoodDef

	flds := []*elorm.FieldDef{gd.Caption, gd.Description}

	flts := []*elorm.Filter{
		elorm.AddFilterEQ(gd.IsDeleted, false),
	}
	orGr := elorm.AddOrGroup()
	flts = append(flts, orGr)
	for _, fld := range flds {
		andGr := elorm.AddAndGroup()
		orGr.Childs = append(orGr.Childs, andGr)
		for _, token := range tokens {
			andGr.Childs = append(andGr.Childs, elorm.AddFilterLIKE(fld, "%"+token+"%"))
		}
	}

	goods, _, err := gd.SelectEntities(flts, nil, 1, 20)
	if err != nil {
		http.Error(w, "Error fetching goods: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(goods)
	if err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
