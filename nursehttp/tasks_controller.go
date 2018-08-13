package nursehttp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/hellupline/hello-nurse/nurseworkers"

	"github.com/hellupline/hello-nurse/booru"
)

type TasksResource struct { // nolint: golint
	TaskManager *nurseworkers.TaskManager
}

func (rs TasksResource) Routes() chi.Router { // nolint: golint
	r := chi.NewRouter()

	r.Route("/booru", func(r chi.Router) {
		r.Post("/fetch-file", rs.BooruFetchFiles)
		r.Post("/fetch-tag", rs.BooruFetchTags)
	})

	return r
}

func (rs TasksResource) BooruFetchTags(w http.ResponseWriter, r *http.Request) { // nolint: golint
	data := BooruFetchTagRequest{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	client, ok := booru.GetClient(data.BooruFetchTask.Domain)
	if !ok {
		_, _ = w.Write([]byte(fmt.Sprintf("%s API not registered", data.BooruFetchTask.Domain)))
		return
	}
	rs.TaskManager.BooruGetTagPage(client("/home/hellupline/.booru/cache/", 100), data.BooruFetchTask.Tag, 0)

	w.WriteHeader(http.StatusCreated)
}

func (rs TasksResource) BooruFetchFiles(w http.ResponseWriter, r *http.Request) { // nolint: golint
	data := BooruFetchFileRequest{}
	if err := render.Bind(r, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rs.TaskManager.BooruGetFile("/home/hellupline/.booru/files/", data.BooruFetchFile.Type, data.BooruFetchFile.Key)

	w.WriteHeader(http.StatusCreated)
}
