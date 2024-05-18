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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)


type Server struct {
	Router *chi.Mux
	DB     *sqlx.DB
	// Db, config can be added here
}

func CreateNewServer(db *sqlx.DB) *Server {
	s := &Server{
		DB:     db,
		Router: chi.NewRouter(),
	}

	return s
}

func (s *Server) Doc() {
	fmt.Println(docgen.MarkdownRoutesDoc(s.Router, docgen.MarkdownOpts{
		ProjectPath: "github.com/go-chi/chi/v5",
	}))
}

func (s *Server) MountHandlers(key string) {
	r := s.Router
	db := s.DB
	// Mount all Middleware here
	r.Use(middleware.Logger) // TODO slog logging?
	r.Use(AuthenticatorSimpleSingleKey(key))
	r.Use(middleware.Recoverer)

	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Mount all handlers here

	r.Route("/", func(r chi.Router) {
		// same like /freezerItems
		r.Get("/", getFreezerItem(db))
		r.Get("/{id}", getFreezerItemByID(db))
		r.Post("/", createFreezerItem(db))
		r.Put("/{id}", updateFreezerItem(db))
		r.Delete("/{id}", deleteFreezerItem(db))
	})

	r.Route("/freezerItems", func(r chi.Router) {
		r.Get("/", getFreezerItem(db))
		r.Get("/{id}", getFreezerItemByID(db))
		r.Post("/", createFreezerItem(db))
		r.Put("/{id}", updateFreezerItem(db))
		r.Delete("/{id}", deleteFreezerItem(db))
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", getCategory(db))
		r.Get("/{name}", getCategoryByName(db))
		r.Post("/", createCategory(db))
		r.Put("/{name}", updateCategory(db))
		r.Delete("/{name}", deleteCategory(db))
	})

	r.Route("/itemTypes", func(r chi.Router) {
		r.Get("/", getItemType(db))
		r.Get("/{name}", getItemTypeByName(db))
		r.Post("/", createItemType(db))
		r.Put("/{name}", updateItemType(db))
		r.Delete("/{name}", deleteItemType(db))
	})
}

func AuthenticatorSimpleSingleKey(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if api key is set to the right value
			if r.Header.Get("x-api-key") != key {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}
