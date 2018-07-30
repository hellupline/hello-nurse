package main

import (
	"net/http"

	"github.com/go-chi/render"
)

type (
	TagListResponse []TagResponse // nolint

	TagResponse struct { // nolint
		*TagCount
	}
)

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
