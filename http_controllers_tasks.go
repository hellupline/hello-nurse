package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gopkg.in/go-playground/pool.v3"
)

type TasksResource struct { // nolint
	Pool pool.Pool
}

func (rs TasksResource) Routes() chi.Router { // nolint
	r := chi.NewRouter()

	r.Route("/booru", func(r chi.Router) {
		r.Post("/fetch-file", rs.BooruFetchFiles)
		r.Post("/fetch-tag", rs.BooruFetchTags)
	})

	return r
}

func (rs TasksResource) BooruFetchFiles(w http.ResponseWriter, r *http.Request) { // nolint
	data := BooruFetchFileRequest{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rs.Pool.Queue(BooruGetFile(rs.Pool, data.BooruFetchFile.Type, data.BooruFetchFile.Key))

	w.WriteHeader(http.StatusCreated)
}

func (rs TasksResource) BooruFetchTags(w http.ResponseWriter, r *http.Request) { // nolint
	data := BooruFetchTagRequest{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rs.Pool.Queue(BooruGetTagPage(rs.Pool, data.BooruFetchTask.Domain, data.BooruFetchTask.Tag, 0))

	w.WriteHeader(http.StatusCreated)
}
