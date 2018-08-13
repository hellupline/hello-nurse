package nursehttp

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/hellupline/hello-nurse/nursedatabase"
)

type (
	PostListResponse []PostResponse // nolint: golint

	PostResponse struct { // nolint: golint
		Post *nursedatabase.Post
	}

	PostRequest struct { // nolint: golint
		Post *nursedatabase.Post
	}
)

func NewPostListResponse(posts []nursedatabase.Post) []render.Renderer { // nolint: golint
	list := []render.Renderer{}

	for i := range posts {
		list = append(list, NewPostResponse(&posts[i]))
	}
	return list
}

func (s PostListResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint: golint
	return nil
}

func NewPostResponse(post *nursedatabase.Post) PostResponse { // nolint: golint
	return PostResponse{Post: post}
}

func (s PostResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint: golint
	return nil
}

func (s PostRequest) Bind(r *http.Request) error { // nolint: golint
	return nil
}
