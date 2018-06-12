package main // import "github.com/hellupline/hello-nurse"

import (
	"net/http"

	nurse "github.com/hellupline/hello-nurse/hellonurse"

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

	datasetGroup := v1Group.Group("/dataset")
	{
		datasetGroup.GET("/download/json", nurse.HttpHandleDownloadDatabaseJSON)
		datasetGroup.POST("/upload/json", nurse.HttpHandleUploadDatabaseJSON)

		datasetGroup.GET("/download/gob", nurse.HttpHandleDownloadDatabase)
		datasetGroup.POST("/upload/gob", nurse.HttpHandleUploadDatabase)
	}

	favoritesGroup := v1Group.Group("/favorites")
	{
		favoritesGroup.GET("/:key", nurse.HttpHandleFavoriteRead)
		favoritesGroup.DELETE("/:key", nurse.HttpHandleFavoriteDelete)
		favoritesGroup.GET("", nurse.HttpHandleFavoriteIndex)
		favoritesGroup.POST("", nurse.HttpHandleFavoriteCreate)
	}

	tagsGroup := v1Group.Group("/tags")
	{
		tagsGroup.GET("", nurse.HttpHandleTagsIndex)
	}

	postsGroup := v1Group.Group("/posts")
	{
		postsGroup.GET("/:key", nurse.HttpHandlePostRead)
		postsGroup.DELETE("/:key", nurse.HttpHandlePostDelete)
		postsGroup.GET("", nurse.HttpHandlePostIndex)
		postsGroup.POST("", nurse.HttpHandlePostCreate)
	}

	router.Run()
}
