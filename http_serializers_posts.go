package main

import (
	"net/http"

	"github.com/go-chi/render"
)

type (
	PostListResponse []PostResponse // nolint

	PostResponse struct { // nolint
		*Post
	}

	PostRequest struct { // nolint
		*Post
	}
)

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
