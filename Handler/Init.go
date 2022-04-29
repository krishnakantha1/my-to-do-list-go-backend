package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CatList []string

type Item struct {
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

type ItemList []Item

type ItemBatch map[string]ItemList

type InitResp struct {
	Categories CatList   `json:"categories"`
	ItemBatch  ItemBatch `json:"itembatch"`
}

func InitializePage(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	dbres, err := db.Query(context.Background(), "select category,line_item,complete from get_init_data()")

	w.Header().Set("content-type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("error in running query"))
		return
	}

	res := InitResp{
		make(CatList, 0),
		make(ItemBatch),
	}

	for dbres.Next() {
		row, err := dbres.Values()

		if err != nil {
			json.NewEncoder(w).Encode(ErrorJson("error in reading rows"))
			return
		}

		//main logic
		cat := row[0].(string)
		_, ok := res.ItemBatch[cat]

		if ok {
			res.ItemBatch[cat] = append(res.ItemBatch[cat], Item{
				Name:     row[1].(string),
				Complete: row[2].(bool),
			})
		} else {
			res.Categories = append(res.Categories, cat)

			if row[2] != nil {
				res.ItemBatch[cat] = append(res.ItemBatch[cat], Item{
					Name:     row[1].(string),
					Complete: row[2].(bool),
				})
			}
		}

	}

	json.NewEncoder(w).Encode(res)
}

func ErrorJson(msg string) map[string]any {
	return map[string]any{
		"error":   true,
		"message": msg,
	}
}

func GenericConfirmation(msg string) map[string]any {
	return map[string]any{
		"error":   false,
		"message": msg,
	}
}
