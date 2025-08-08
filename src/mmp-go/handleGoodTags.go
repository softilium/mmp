package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/softilium/elorm"
)

func initRouterGoodTags(router *http.ServeMux) {

	//// tags

	tagsRestApiConfig := elorm.CreateStdRestApiConfig(
		*DB.TagDef.EntityDef,
		DB.LoadTag,
		DB.TagDef.SelectEntities,
		DB.CreateTag)
	tagsRestApiConfig.BeforeMiddleware = AdminRequiredForEdit
	tagsRestApiConfig.DefaultPageSize = 0
	tagsRestApiConfig.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{
			{Field: DB.TagDef.Color, Asc: true},
			{Field: DB.TagDef.Name, Asc: true},
		}, nil
	}
	router.HandleFunc("/api/tags", elorm.HandleRestApi(tagsRestApiConfig))

	router.HandleFunc("/api/good-tags", handleGoodTags)
	router.HandleFunc("/api/tags-by-shop", handleTagsByShop)
	router.HandleFunc("/api/tags-by-all", handleTagsByAll)
	router.HandleFunc("/api/goods-by-tag", handleGoodsByTag)
}

type TagResultLine struct {
	TagRef   string `json:"tagRef"`
	TagName  string `json:"tagName"`
	TagColor string `json:"tagColor,omitempty"`
	Tagged   bool   `json:"tagged"`
}

func handleGoodTags(w http.ResponseWriter, r *http.Request) {
	gref := r.URL.Query().Get("ref")
	if gref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("good ref is required"))
		return
	}
	good, err := DB.LoadGood(gref)
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error loading good: %v", err))
		return
	}
	if r.Method == http.MethodGet {
		rows, err := DB.Query(fmt.Sprintf(
			`
			select 
				t.ref as tag, t.name as tagname, case when gt.good is null then '0' else '1' end as Tagged, t.color
			from tags t
			left join goodtags gt on (t.ref=gt.tag and gt.good='%s')
			order by t.color, t.name
			`,
			gref))
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error querying tags: %v", err))
			return
		}
		result := make([]TagResultLine, 0)
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var tagRef string
			var tagName string
			var tagged string
			var tagColor string
			err := rows.Scan(&tagRef, &tagName, &tagged, &tagColor)
			if err != nil {
				HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error scanning row: %v", err))
				return
			}
			result = append(result, TagResultLine{
				TagRef:   tagRef,
				TagName:  tagName,
				Tagged:   tagged == "1",
				TagColor: tagColor,
			})
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error encoding response: %v", err))
			return
		}
	}
	if r.Method == http.MethodPost {
		result := make([]TagResultLine, 0)
		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			HandleErr(w, http.StatusBadRequest, fmt.Errorf("error decoding request body: %v", err))
			return
		}
		tx, err := DB.BeginTran()
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error starting transaction: %v", err))
			return
		}
		defer func() { _ = DB.RollbackTran(tx) }()

		old, _, err := DB.GoodTagDef.SelectEntities(
			[]*elorm.Filter{elorm.AddFilterEQ(DB.GoodTagDef.Good, good)}, nil, 0, 0)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error selecting old good tags: %v", err))
			return
		}
		for _, ot := range old {
			err = DB.DeleteEntity(r.Context(), ot.RefString())
			if err != nil {
				HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error deleting old good tag: %v", err))
				return
			}
		}
		for _, line := range result {
			if line.Tagged {
				gt, err := DB.CreateGoodTag()
				if err != nil {
					HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error creating good tag: %v", err))
					return
				}
				gt.SetGood(good)

				tg, err := DB.LoadTag(line.TagRef)
				if err != nil {
					HandleErr(w, http.StatusBadRequest, fmt.Errorf("tag %s not found", line.TagRef))
					return
				}

				gt.SetTag(tg)
				err = gt.Save(r.Context())
				if err != nil {
					HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error saving good tag: %v", err))
					return
				}
			}
		}
		err = DB.CommitTran(tx)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error committing transaction: %v", err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleTagsByShop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleErr(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed, only GET is supported"))
		return
	}
	sref := r.URL.Query().Get("ref")
	if sref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("shop ref is required"))
		return
	}
	result, err := getTags(sref)
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error getting tags: %v", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error encoding response: %v", err))
		return
	}
}

func handleTagsByAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleErr(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed, only GET is supported"))
		return
	}
	result, err := getTags("")
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error getting tags: %v", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error encoding response: %v", err))
		return
	}
}

func getTags(sref string) ([]TagResultLine, error) {

	b := "false"
	if sref == "" {
		b = "true"
	}

	rows, err := DB.Query(fmt.Sprintf(`

		select t.ref, t.name, t.color, count(*) as cnt
		from goodtags gt
		left join tags t on gt.tag=t.ref
		inner join goods g on gt.good=g.ref and (g.ownershop='%s' or %s)
		group by t.ref
		having count(*) > 0
		order by cnt DESC

	`, sref, b))
	if err != nil {
		return nil, fmt.Errorf("error querying tags by shop: %v", err)
	}
	defer func() { _ = rows.Close() }()
	result := make([]TagResultLine, 0)
	for rows.Next() {
		var tagRef string
		var tagName string
		var tagColor string
		var count int
		err := rows.Scan(&tagRef, &tagName, &tagColor, &count)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		result = append(result, TagResultLine{
			TagRef:   tagRef,
			TagName:  tagName,
			TagColor: tagColor,
		})
	}

	return result, nil
}

func handleGoodsByTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleErr(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed, only GET is supported"))
		return
	}
	tref := r.URL.Query().Get("ref")
	if tref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("tag ref is required"))
		return
	}
	rows, err := DB.Query(fmt.Sprintf(`
		select gt.good
		from goodtags gt
		where gt.tag='%s'
		order by gt.ref
	`, tref))
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error querying goods by tag: %v", err))
		return
	}
	defer func() { _ = rows.Close() }()

	result := make([]*Good, 0)
	for rows.Next() {
		var ref string
		err := rows.Scan(&ref)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error scanning row: %v", err))
			return
		}
		good, err := DB.LoadGood(ref)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error loading good: %v", err))
			return
		}
		result = append(result, good)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		HandleErr(w, http.StatusInternalServerError, fmt.Errorf("error encoding response: %v", err))
		return
	}
}
