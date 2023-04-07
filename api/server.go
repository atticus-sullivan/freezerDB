package api

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
