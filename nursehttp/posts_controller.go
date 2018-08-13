package nursehttp

import (
	"net/http"
	"sort"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/hellupline/hello-nurse/nursedatabase"
)

type PostsResource struct { // nolint: golint
	Database *nursedatabase.Database
}

func (rs PostsResource) Routes() chi.Router { // nolint: golint
	r := chi.NewRouter()

	r.Route("/{type}/{key}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Delete("/", rs.Delete)
	})

	r.Get("/", rs.Index)
	r.Post("/", rs.Create)

	return r
}

func (rs PostsResource) Index(w http.ResponseWriter, r *http.Request) { // nolint: golint
	posts := rs.Database.PostIndex(r.URL.Query().Get("q"))

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Key < posts[j].Key
	})

	w.WriteHeader(http.StatusOK)
	_ = render.RenderList(w, r, NewPostListResponse(posts))
}

func (rs PostsResource) Create(w http.ResponseWriter, r *http.Request) { // nolint: golint
	data := PostRequest{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rs.Database.PostCreate(*data.Post)

	w.WriteHeader(http.StatusCreated)
	_ = render.Render(w, r, NewPostResponse(data.Post))
}

func (rs PostsResource) Get(w http.ResponseWriter, r *http.Request) { // nolint: golint
	key := nursedatabase.PostKey{
		Type: chi.URLParam(r, "type"),
		Key:  chi.URLParam(r, "key"),
	}
	post, ok := rs.Database.PostRead(key)
	if !ok {
		render.Status(r, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = render.Render(w, r, NewPostResponse(&post))
}

func (rs PostsResource) Delete(w http.ResponseWriter, r *http.Request) { // nolint: golint
	key := nursedatabase.PostKey{
		Type: chi.URLParam(r, "type"),
		Key:  chi.URLParam(r, "key"),
	}
	rs.Database.PostDelete(key)

	w.WriteHeader(http.StatusNoContent)
}
