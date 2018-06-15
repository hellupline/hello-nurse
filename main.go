package main // import "github.com/hellupline/hello-nurse"

import (
	"net/http"

	"github.com/hellupline/hello-nurse/nursetags"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	v1Group := router.Group("/v1")

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
		postsGroup.GET("/:key", nursetags.HttpHandlePostRead)
		postsGroup.DELETE("/:key", nursetags.HttpHandlePostDelete)
		postsGroup.GET("", nursetags.HttpHandlePostIndex)
		postsGroup.POST("", nursetags.HttpHandlePostCreate)
	}

	datasetGroup := v1Group.Group("/dataset")
	{
		datasetGroup.GET("/download/json", nursetags.HttpHandleDownloadDatabaseJSON)
		datasetGroup.POST("/upload/json", nursetags.HttpHandleUploadDatabaseJSON)

		datasetGroup.GET("/download/gob", nursetags.HttpHandleDownloadDatabaseGOB)
		datasetGroup.POST("/upload/gob", nursetags.HttpHandleUploadDatabaseGOB)
	}

	router.Run()
}
