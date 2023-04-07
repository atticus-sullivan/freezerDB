package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/atticus-sullivan/freezerDB/db/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Error:", e.Err)
	e.Err = nil // don't expose error
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
func ErrInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error.",
		ErrorText:      "",
	}
}

var ErrNotFound = &ErrResponse{
	HTTPStatusCode: http.StatusNotFound,
	StatusText:     "Resource not found.",
}

// Servers may also send 404 Not found instead of 403 Forbidden to hide the
// existence of a resource from an unauthorized client.
var ErrForbidden = ErrNotFound
var ErrItemNotFount error = errors.New("Item not found")
var ErrInvalidID error = errors.New("Invalid Item ID")

type CategoryRequest struct{ models.Category }

func (c *CategoryRequest) Bind(r *http.Request) error {
	return nil
}

type CategoryResponse struct{ models.Category }

func (c CategoryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

type ItemTypeRequest struct{ models.ItemType }

func (c *ItemTypeRequest) Bind(r *http.Request) error {
	return nil
}

type ItemTypeResponse struct{ models.ItemType }

func (c ItemTypeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

type FreezerItemRequest struct{ models.FreezerItem }

func (c *FreezerItemRequest) Bind(r *http.Request) error {
	return nil
}

type FreezerItemResponse struct{ models.FreezerItem }

func (c FreezerItemResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

func NewListResponse[T render.Renderer](items []T) []render.Renderer {
	list := []render.Renderer{}
	for _, item := range items {
		list = append(list, item)
	}
	return list
}

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
