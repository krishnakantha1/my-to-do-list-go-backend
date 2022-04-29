package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type RequestItem struct {
	Cat   string `json:"category"`
	Entry string `json:"line_item"`
}

func parseItemReq(R *RequestItem, r *http.Request) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(R)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ItemMarkAsDone(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

	req := RequestItem{}

	err := parseItemReq(&req, r)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("error in parsing request"))
		return
	}

	dbres, err := db.Query(context.Background(), "CALL toggle_todo_done($1,$2,'')", req.Cat, req.Entry)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("server error"))
		return
	}

	for dbres.Next() {
		row, err := dbres.Values()

		if err != nil {
			json.NewEncoder(w).Encode(ErrorJson("server error"))
			return
		}

		if len(row) > 0 && row[0] != nil {
			json.NewEncoder(w).Encode(ErrorJson(row[0].(string)))
			return
		}

		returnAllEntries(req.Cat, db, w)
	}
}

func DeleteItem(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	req := RequestItem{}

	err := parseItemReq(&req, r)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("error in parsing request"))
		return
	}

	dbres, err := db.Query(context.Background(), "CALL delete_line_item($1,$2,'')", req.Cat, req.Entry)

	if err != nil {
		json.NewEncoder(w).Encode(ErrorJson("server error"))
		return
	}

	for dbres.Next() {
		row, err := dbres.Values()

		if err != nil {
			json.NewEncoder(w).Encode(ErrorJson("server error"))
			return
		}

		if len(row) > 0 && row[0] != nil {
			json.NewEncoder(w).Encode(ErrorJson(row[0].(string)))
			return
		}

		returnAllEntries(req.Cat, db, w)
	}

}
