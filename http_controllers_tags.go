package main

import (
	"net/http"
	"sort"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type TagsResource struct{} // nolint

func (rs TagsResource) Routes() chi.Router { // nolint
	r := chi.NewRouter()

	r.Route("/{key}", func(r chi.Router) {
		r.Get("/", rs.Get)
	})

	r.Get("/", rs.Index)

	return r
}

func (rs TagsResource) Index(w http.ResponseWriter, r *http.Request) { // nolint
	tagNames := database.TagIndex()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i].Count > tagNames[j].Count
	})

	w.WriteHeader(http.StatusOK)
	_ = render.RenderList(w, r, NewTagListResponse(tagNames))
}

func (rs TagsResource) Get(w http.ResponseWriter, r *http.Request) { // nolint
	_, _ = w.Write([]byte("todo get"))
}
