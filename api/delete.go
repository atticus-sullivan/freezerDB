package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func deleteFreezerItem(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract item ID from request URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(ErrInvalidID))
			return
		}

		// Delete item from database
		result, err := db.Exec("DELETE FROM freezer_items WHERE id = ?", id)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Check if item was found and deleted
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		if rowsAffected == 0 {
			render.Render(w, r, ErrInvalidRequest(ErrInvalidID))
			return
		}

		// Send response
		render.NoContent(w, r)
		return
	}
}
func deleteItemType(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract type name from request URL
		name := chi.URLParam(r, "name")

		// Delete item from database
		result, err := db.Exec("DELETE FROM item_types WHERE name = ?", name)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Check if item was found and deleted
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		if rowsAffected == 0 {
			render.Render(w, r, ErrInvalidRequest(ErrInvalidID))
			return
		}

		// Send response
		render.NoContent(w, r)
		return
	}
}
func deleteCategory(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract category name from request URL
		name := chi.URLParam(r, "name")

		// Delete category from database
		result, err := db.Exec("DELETE FROM categories WHERE name=?", name)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// Check if item was found and deleted
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		if rowsAffected == 0 {
			render.Render(w, r, ErrInvalidRequest(ErrInvalidID))
			return
		}

		// Send response
		render.NoContent(w, r)
		return
	}
}
