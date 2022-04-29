package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type S struct {
	Category string   `json:"category"`
	Items    []string `json:"items"`
}

//handler to insert new items to the required category
func AddToDoItemsToDB(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	req := S{}
	d.Decode(&req)

	res, err := db.Query(context.Background(), "call add_items_for_category($1,$2,'')", req.Category, req.Items)
	if err != nil {
		panic(err)
	}

	w.Header().Set("content-type", "application/json")

	for res.Next() {
		row, err := res.Values()

		if err != nil {
			json.NewEncoder(w).Encode(ErrorJson("Server Error"))
			return
		}

		if len(row) > 0 && row[0] != nil {
			json.NewEncoder(w).Encode(ErrorJson(row[0].(string)))
			return
		}

		returnAllEntries(req.Category, db, w)
	}
}

//return all the items in the current category
func returnAllEntries(cat string, db *pgxpool.Pool, w http.ResponseWriter) {
	dbres, err := db.Query(context.Background(), "select category,line_item,complete from get_items_for_category($1)", cat)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("Server Error"))
		return
	}

	res := make(map[string]ItemList)

	for dbres.Next() {
		row, e := dbres.Values()

		if e != nil {
			json.NewEncoder(w).Encode(ErrorJson("Server Error"))
			return
		}

		res[cat] = append(res[cat], Item{
			Name:     row[1].(string),
			Complete: row[2].(bool),
		})
	}

	json.NewEncoder(w).Encode(res)
}
