package nursehttp

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/hellupline/hello-nurse/nursedatabase"
)

type (
	TagListResponse []TagResponse // nolint: golint

	TagResponse struct { // nolint: golint
		TagCount *nursedatabase.TagCount
	}
)

func NewTagListResponse(tags []nursedatabase.TagCount) []render.Renderer { // nolint: golint
	list := []render.Renderer{}

	for i := range tags {
		list = append(list, NewTagResponse(&tags[i]))
	}
	return list
}

func (s TagListResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint: golint
	return nil
}

func NewTagResponse(tag *nursedatabase.TagCount) TagResponse { // nolint: golint
	return TagResponse{tag}
}

func (s TagResponse) Render(w http.ResponseWriter, r *http.Request) error { // nolint: golint
	return nil
}
