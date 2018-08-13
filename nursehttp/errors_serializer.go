package nursehttp

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."} // nolint: golint
)

type (
	ErrResponse struct { // nolint: golint
		Err            error `json:"-"` // low-level runtime error
		HTTPStatusCode int   `json:"-"` // http response status code

		StatusText string `json:"status"`          // user-level status message
		AppCode    int64  `json:"code,omitempty"`  // application-specific error code
		ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
	}
)

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint: golint
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer { // nolint: golint
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer { // nolint: golint
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
