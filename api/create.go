package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func createFreezerItem(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var item FreezerItemRequest
		if err := render.Bind(r, &item); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		fmt.Printf("%+v\n", item)

		// Insert into database
		_, err := db.NamedExec("INSERT INTO freezer_items (date, identifier, amount, misc, item_name) VALUES (:date, :identifier, :amount, :misc, :item_name)", item)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		return
	}
}
func createItemType(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var item ItemTypeRequest
		if err := render.Bind(r, &item); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		fmt.Printf("%+v\n", item)

		// Insert into database
		_, err := db.NamedExec("INSERT INTO item_types (name, category_name) VALUES (:name, :category_name)", item)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
	}
}
func createCategory(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var category CategoryRequest
		if err := render.Bind(r, &category); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// Insert into database
		_, err := db.NamedExec("INSERT INTO categories (name) VALUES (:name)", category)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
	}
}
