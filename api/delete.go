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
