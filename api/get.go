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

func getFreezerItem(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query from database
		var items []FreezerItemResponse
		err := db.Select(&items, "SELECT * FROM freezer_items ORDER BY item_name")
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
