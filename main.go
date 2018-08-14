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
	"github.com/go-chi/cors"
	"gopkg.in/go-playground/pool.v3"

	"github.com/hellupline/hello-nurse/nursedatabase"
	"github.com/hellupline/hello-nurse/nursehttp"
	"github.com/hellupline/hello-nurse/nurseworkers"
)

var (
	quit = make(chan os.Signal)

	baseDir string
)

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	baseDir = filepath.Join(u.HomeDir, ".booru")
}

func main() {
	db := nursedatabase.NewDatabase()
	if err := readDatabase(db); err != nil {
		log.Fatal("error on opening database: ", err)
	}
	// defer after load, do not save if could not open first
	defer writeDatabase(db) // nolint: errcheck

	p := pool.NewLimited(8)
	defer p.Close()

	taskManager := nurseworkers.NewTaskManager(db, p)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/posts", nursehttp.PostsResource{Database: db}.Routes())
			r.Mount("/tags", nursehttp.TagsResource{Database: db}.Routes())
			r.Mount("/tasks", nursehttp.TasksResource{TaskManager: taskManager}.Routes())
		})
	})

	r.Get("/", httpHandleIndex)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	signal.Notify(quit, os.Interrupt)
	<-quit
}

func httpHandleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	_, _ = w.Write(data) // nolint: gosec
}

func readDatabase(db *nursedatabase.Database) error {
	f, err := os.Open(filepath.Join(baseDir, "db.gob"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return db.Read(f)
}

func writeDatabase(db *nursedatabase.Database) error {
	f, err := os.Create(filepath.Join(baseDir, "db.gob"))
	if err != nil {
		return err
	}
	return db.Write(f)
}
