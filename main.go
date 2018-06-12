package main // import "github.com/hellupline/hello-nurse"

import (
	"net/http"

	nurse "github.com/hellupline/hello-nurse/hellonurse"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
)

func init() {
	router := gin.Default()
	router.Use(cors.Default())
	http.Handle("/", router)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	v1Group := router.Group("/v1")

	favoritesGroup := v1Group.Group("/favorites")
	{
		favoritesGroup.GET("/:key", nurse.HttpHandleFavoriteRead)
		favoritesGroup.DELETE("/:key", nurse.HttpHandleFavoriteDelete)
		favoritesGroup.GET("", nurse.HttpHandleFavoriteIndex)
		favoritesGroup.POST("", nurse.HttpHandleFavoriteCreate)
	}

	postsGroup := v1Group.Group("/posts")
	{
		postsGroup.GET("/:key", nurse.HttpHandlePostRead)
		postsGroup.DELETE("/:key", nurse.HttpHandlePostDelete)
		postsGroup.GET("", nurse.HttpHandlePostIndex)
		postsGroup.POST("", nurse.HttpHandlePostCreate)
	}

	v1Group.GET("/tags", nurse.HttpHandleTagsIndex)

	v1Group.GET("/download", nurse.HttpHandleDownloadDatabase)
	v1Group.POST("/upload", nurse.HttpHandleUploadDatabase)
}

func main() {
	appengine.Main()
}
