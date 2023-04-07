package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/atticus-sullivan/freezerDB/db/models"
	"github.com/go-chi/render"
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
