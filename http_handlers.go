package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gopkg.in/go-playground/pool.v3"
)

type (
	PostsResource struct{} // nolint

	TagsResource struct{} // nolint

	TasksResource struct { // nolint
		Pool pool.Pool
	}
)

func httpHandleIndex(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("./index.html")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	_, _ = w.Write(data)
}

func (rs PostsResource) Routes() chi.Router { // nolint
	r := chi.NewRouter()

	r.Route("/{type}/{key}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Delete("/", rs.Delete)
	})

	r.Get("/", rs.Index)
	r.Post("/", rs.Create)

	return r
}

func (rs PostsResource) Index(w http.ResponseWriter, r *http.Request) { // nolint
	posts := database.PostIndex(r.URL.Query().Get("q"))

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Key < posts[j].Key
	})

	w.WriteHeader(http.StatusOK)
	_ = render.RenderList(w, r, NewPostListResponse(posts))
}

func (rs PostsResource) Create(w http.ResponseWriter, r *http.Request) { // nolint
	data := PostRequest{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	database.PostCreate(*data.Post)

	w.WriteHeader(http.StatusCreated)
	_ = render.Render(w, r, NewPostResponse(data.Post))
}

func (rs PostsResource) Get(w http.ResponseWriter, r *http.Request) { // nolint
	key := PostKey{
		Type: chi.URLParam(r, "type"),
		Key:  chi.URLParam(r, "key"),
	}
	post, ok := database.PostRead(key)
	if !ok {
		render.Status(r, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = render.Render(w, r, NewPostResponse(&post))
}

func (rs PostsResource) Delete(w http.ResponseWriter, r *http.Request) { // nolint
	key := PostKey{
		Type: chi.URLParam(r, "type"),
		Key:  chi.URLParam(r, "key"),
	}
	database.PostDelete(key)

	w.WriteHeader(http.StatusNoContent)
}

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
	for _, c := range tagNames {
		log.Println(c)
	}
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i].Count > tagNames[j].Count
	})

	w.WriteHeader(http.StatusOK)
	_ = render.RenderList(w, r, NewTagListResponse(tagNames))
}

func (rs TagsResource) Get(w http.ResponseWriter, r *http.Request) { // nolint
	_, _ = w.Write([]byte("todo get"))
}

func (rs TasksResource) Routes() chi.Router { // nolint
	r := chi.NewRouter()

	r.Route("/booru", func(r chi.Router) {
		r.Post("/fetch-tasks", rs.BooruFetchTags)
	})

	return r
}

func (rs TasksResource) BooruFetchTags(w http.ResponseWriter, r *http.Request) { // nolint
	rs.Pool.Queue(GetBooruTagPage(rs.Pool, "konachan.net", "eureka_seven", 0))

	w.WriteHeader(http.StatusCreated)
}
