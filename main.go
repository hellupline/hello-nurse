package main // import "github.com/hellupline/hello-nurse"

import (
	"net/http"

	"github.com/hellupline/hello-nurse/nursetags"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
)

func init() {
	router := gin.Default()
	router.Use(cors.Default())

	v1Group := router.Group("/v1")

	router.GET("/_ah/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	favoritesGroup := v1Group.Group("/favorites")
	{
		favoritesGroup.GET("/:key", nursetags.HttpHandleFavoriteRead)
		favoritesGroup.DELETE("/:key", nursetags.HttpHandleFavoriteDelete)
		favoritesGroup.GET("", nursetags.HttpHandleFavoriteIndex)
		favoritesGroup.POST("", nursetags.HttpHandleFavoriteCreate)
	}

	tagsGroup := v1Group.Group("/tags")
	{
		tagsGroup.GET("", nursetags.HttpHandleTagsIndex)
	}

	postsGroup := v1Group.Group("/posts")
	{
		postsGroup.GET("/:namespace/:key", nursetags.HttpHandlePostRead)
		postsGroup.DELETE("/:namespace/:key", nursetags.HttpHandlePostDelete)
		postsGroup.GET("", nursetags.HttpHandlePostIndex)
		postsGroup.POST("", nursetags.HttpHandlePostCreate)
	}

	datasetGroup := v1Group.Group("/dataset")
	{
		datasetGroup.GET("/download/gob", nursetags.HttpHandleDownloadDatasetGOB)
		datasetGroup.POST("/upload/gob", nursetags.HttpHandleUploadDatasetGOB)
	}

	// _ = router.Run()
	http.Handle("/", router)
}

func main() {
	appengine.Main()
}
