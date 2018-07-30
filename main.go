package main // import "github.com/hellupline/hello-nurse"

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gopkg.in/go-playground/pool.v3"
)

var (
	quit    = make(chan os.Signal)
	baseDir string

	database = NewDatabase()
)

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	baseDir = filepath.Join(u.HomeDir, ".booru")
}

func main() {
	if f, err := os.Open(filepath.Join(baseDir, "db.gob")); err == nil {
		if err := database.Read(f); err != nil {
			log.Fatal(err)
		}
	}

	defer func() {
		// defer after load, do not save if could not open first
		if f, err := os.Create(filepath.Join(baseDir, "db.gob")); err == nil {
			if err := database.Write(f); err != nil {
				log.Fatal(err)
			}
		}
	}()

	p := pool.NewLimited(8)
	defer p.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/posts", PostsResource{}.Routes())
		r.Mount("/tags", TagsResource{}.Routes())
		r.Mount("/tasks", TasksResource{Pool: p}.Routes())
	})

	r.Get("/", httpHandleIndex)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	signal.Notify(quit, os.Interrupt)
	<-quit
}

func httpHandleIndex(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("./index.html")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	_, _ = w.Write(data)
}
