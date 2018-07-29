package main

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
)

type (
	ErrResponse struct {
		Err            error `json:"-"` // low-level runtime error
		HTTPStatusCode int   `json:"-"` // http response status code

		StatusText string `json:"status"`          // user-level status message
		AppCode    int64  `json:"code,omitempty"`  // application-specific error code
		ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
	}

	PostListResponse []PostResponse // nolint

	PostResponse struct { // nolint
		*Post
	}

	PostRequest struct { // nolint
		*Post
	}

	TagListResponse []TagResponse // nolint

	TagResponse struct { // nolint
		*TagCount
	}
)

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func NewPostListResponse(posts []Post) []render.Renderer {
	list := []render.Renderer{}

	for i := range posts {
		list = append(list, NewPostResponse(&posts[i]))
	}
	return list
}

func (s PostListResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint
	return nil
}

func NewPostResponse(post *Post) PostResponse {
	return PostResponse{post}
}

func (s PostResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint
	return nil
}

func (s PostRequest) Bind(r *http.Request) error {
	return nil
}

func NewTagListResponse(tags []TagCount) []render.Renderer {
	list := []render.Renderer{}

	for i := range tags {
		list = append(list, NewTagResponse(&tags[i]))
	}
	return list
}

func (s TagListResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint
	return nil
}

func NewTagResponse(tag *TagCount) TagResponse {
	return TagResponse{tag}
}

func (s TagResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint
	return nil
}
