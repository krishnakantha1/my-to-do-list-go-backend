package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Category struct {
	Cat string `json:"category"`
}

//insert the category that the user have provided to DB
func CreateCategory(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	req := Category{}

	d.Decode(&req)

	w.Header().Set("content-type", "application/json")

	if len(req.Cat) == 0 {
		json.NewEncoder(w).Encode(ErrorJson("category is not present in request"))
		return
	}

	dbres, err := db.Query(context.Background(), "call add_category($1,'')", req.Cat)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("error in running query"))
		return
	}

	for dbres.Next() {
		row, err := dbres.Values()

		if err != nil {
			json.NewEncoder(w).Encode(ErrorJson("server error"))
			return
		}

		if row[0] != nil {
			json.NewEncoder(w).Encode(ErrorJson(row[0].(string)))
			return
		}

		json.NewEncoder(w).Encode(GenericConfirmation("Category created"))
	}
}

//Get the number of items listed under the given category in the DB
func GetCategoryEntryCount(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	req := Category{}

	err := d.Decode(&req)

	w.Header().Set("content-type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("category is not present in request"))
		return
	}

	count := 0
	e := db.QueryRow(context.Background(), "select get_category_entry_count($1) as count", req.Cat).Scan(&count)

	if e != nil {
		json.NewEncoder(w).Encode(ErrorJson("error running the query"))
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"error": false,
		"count": count,
	})
}

func DeleteCategory(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	req := Category{}

	err := d.Decode(&req)

	w.Header().Set("content-type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("category is not present in request"))
		return
	}

	dbres, e := db.Query(context.Background(), "CALL delete_category($1,0,'')", req.Cat)

	if e != nil {
		json.NewEncoder(w).Encode(ErrorJson("error in running query"))
		return
	}

	for dbres.Next() {
		row, e := dbres.Values()

		if e != nil {
			json.NewEncoder(w).Encode(ErrorJson("error in running query"))
			return
		}

		if row[1] != nil {
			json.NewEncoder(w).Encode(ErrorJson(row[1].(string)))
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"error": false,
			"count": row[0].(int32),
		})
	}
}
