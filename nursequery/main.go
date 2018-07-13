package main // import "github.com/hellupline/hello-nurse/nursequery"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"

	"github.com/hellupline/hello-nurse/nursetags"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type (
	NurseDownloadTask struct {
		Domain string `json:"domain" binding:"required"`
		URL    string `json:"url" binding:"required"`
	}

	NurseFetchTask struct {
		Domain string `json:"domain" binding:"required"`
		Tag    string `json:"tag" binding:"required"`
	}
)

var (
	booruDownloadStage0 = make(chan *NurseDownloadTask, 100)
	booruFetchStage0    = make(chan *TagPage, 100)

	quit = make(chan os.Signal)

	database = nursetags.NewDatabase()

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
	{
		if f, err := os.Open(filepath.Join(baseDir, "db.gob")); err == nil {
			database.Read(f)
		}
	}
	r := gin.Default()

	r.Use(static.Serve("/files", static.LocalFile(filepath.Join(baseDir, "files"), true)))
	r.Use(cors.Default())

	v1Group := r.Group("/v1")
	{
		postsGroup := v1Group.Group("/posts")
		{
			postsGroup.GET("/:namespace/:key", HttpHandlePostRead)
			postsGroup.DELETE("/:namespace/:key", HttpHandlePostDelete)
			postsGroup.GET("", HttpHandlePostIndex)
			postsGroup.POST("", HttpHandlePostCreate)
		}

		tagsGroup := v1Group.Group("/tags")
		{
			tagsGroup.GET("", HttpHandleTagsIndex)
		}

		fetchGroup := v1Group.Group("/tasks")
		{
			fetchGroup.POST("/nurse-download", HttpHandleBooruDownload)
			fetchGroup.POST("/nurse-fetch", HttpHandleBooruFetch)
		}
	}

	r.GET("/", func(c *gin.Context) {
		data, _ := ioutil.ReadFile("./index.html")
		c.String(http.StatusOK, string(data))
	})
	r.GET("/_ah/health", HttpHandleHealthCheck)

	go func() {
		stage1 := BooruDownloadPipeline(booruDownloadStage0, DownloadFile)

		for range stage1 {
			// XXX: Log completes
		}
	}()

	go func() {
		stage1 := BooruFetchPipeline(booruFetchStage0, FetchFirstPageStage)
		stage2 := BooruFetchPipeline(stage1, FetchAllPagesStage)
		stage3 := BooruFetchPipeline(stage2, SaveToQueryServer)

		for range stage3 {
			// XXX: Log completes
		}
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	signal.Notify(quit, os.Interrupt)
	<-quit

	{
		if f, err := os.Create(filepath.Join(baseDir, "db.gob")); err == nil {
			database.Write(f)
		}
	}
}

func bindErrorResponse(err error) map[string][]string {
	errors := make(map[string][]string)
	switch msg := err.(type) {
	case validator.ValidationErrors:
		for _, v := range msg {
			errors[v.Field()] = append(errors[v.Field()], v.Tag())
		}
	default:
		errors["unknown"] = []string{fmt.Sprintln(err)}
	}
	return errors
}

func HttpHandleHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func HttpHandleBooruFetch(c *gin.Context) {
	task := NurseFetchTask{}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	booruFetchStage0 <- NewTagPage(task.Domain, task.Tag)

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func HttpHandleBooruDownload(c *gin.Context) {
	task := NurseDownloadTask{}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	booruDownloadStage0 <- &task

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}
