package api

// Copyright (c) 2023, Lukas Heindl
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"net/http"
	"strconv"

	"github.com/atticus-sullivan/freezerDB/db/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

func updateFreezerItem(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract item ID from request URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(ErrInvalidID))
			return
		}

		// Parse request body
		var item FreezerItemRequest
		if err := render.Bind(r, &item); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// query current item
		var def models.FreezerItem
		if err := db.Get(&def, "SELECT * FROM freezer_items WHERE id = ?", id); err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// use current item as default for zero values
		item.FillDefaults(&def)

		// Update item type in database
		result, err := db.Exec("UPDATE freezer_items SET id = ?, identifier = ?, amount = ?, misc = ?, item_name = ? WHERE id = ?", item.ID, item.Identifier, item.Amount, item.Misc, item.ItemName, id)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// check if any update was done
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
func updateItemType(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract item type name from request URL
		name := chi.URLParam(r, "name")

		// Parse request body
		var item ItemTypeRequest
		if err := render.Bind(r, &item); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// query current item
		var def models.ItemType
		if err := db.Get(&def, "SELECT * FROM category_name WHERE name = ?", name); err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// use current item as default for zero values
		item.FillDefaults(&def)

		// Update item type in database
		result, err := db.Exec("UPDATE item_types SET name = ?, category_name = ? WHERE name = ?", item.Name, item.CategoryName, name)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// check if any update was done
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
func updateCategory(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		// Parse request body
		var category CategoryRequest
		if err := render.Bind(r, &category); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		// query current item
		var def models.Category
		if err := db.Get(&def, "SELECT * FROM categories WHERE name = ?", name); err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}
		// use current item as default for zero values
		category.FillDefaults(&def)

		result, err := db.Exec("UPDATE categories SET name=? WHERE name=?", category.Name, name)
		if err != nil {
			render.Render(w, r, ErrInternal(err))
			return
		}

		// check if any update was done
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
