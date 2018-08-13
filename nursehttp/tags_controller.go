package nursehttp

import (
	"net/http"
	"sort"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/hellupline/hello-nurse/nursedatabase"
)

type TagsResource struct { // nolint: golint
	Database *nursedatabase.Database
}

func (rs TagsResource) Routes() chi.Router { // nolint: golint
	r := chi.NewRouter()

	r.Route("/{key}", func(r chi.Router) {
		r.Get("/", rs.Get)
	})

	r.Get("/", rs.Index)

	return r
}

func (rs TagsResource) Index(w http.ResponseWriter, r *http.Request) { // nolint: golint
	tagNames := rs.Database.TagIndex()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i].Count > tagNames[j].Count
	})

	w.WriteHeader(http.StatusOK)
	_ = render.RenderList(w, r, NewTagListResponse(tagNames))
}

func (rs TagsResource) Get(w http.ResponseWriter, r *http.Request) { // nolint: golint
	_, _ = w.Write([]byte("todo get"))
}
