package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func getFreezerItem(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query from database
		var items []FreezerItemResponse
		err := db.Select(&items, "SELECT * FROM freezer_items ORDER BY name")
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Send respons	e
		render.RenderList(w, r, NewListResponse(items))
		return
	}
}
func getItemType(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query from database
		var items []ItemTypeResponse
		err := db.Select(&items, "SELECT * FROM item_types ORDER BY name")
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// Send response
		render.RenderList(w, r, NewListResponse(items))
		return
	}
}
func getCategory(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query from database
		var items []CategoryResponse
		err := db.Select(&items, "SELECT * FROM categories ORDER BY name")
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// Send response
		render.RenderList(w, r, NewListResponse(items))
		return
	}
}

func getFreezerItemByID(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract item ID from request URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(ErrInvalidID))
			return
		}

		// Query from database
		var item FreezerItemResponse
		err = db.Get(&item, "SELECT * FROM freezer_items WHERE id = ?", id)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// Send response
		render.Render(w, r, &item)
		return
	}
}
func getItemTypeByName(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract item name from request URL
		name := chi.URLParam(r, "name")

		// Query from database
		var item ItemTypeResponse
		err := db.Get(&item, "SELECT * FROM item_types WHERE name = ?", name)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// Send response
		render.Render(w, r, &item)
		return
	}
}
func getCategoryByName(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract item name from request URL
		name := chi.URLParam(r, "name")

		// Query from database
		var item CategoryResponse
		err := db.Get(&item, "SELECT * FROM categories WHERE name = ?", name)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// Send response
		render.Render(w, r, &item)
		return
	}
}
