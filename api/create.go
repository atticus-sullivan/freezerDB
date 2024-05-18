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
