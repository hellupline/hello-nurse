package main // import "github.com/hellupline/hello-nurse/hello-nurse"

import (
	"fmt"
	"net/http"

	"github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	"google.golang.org/appengine"
)

type (
	FavoritesDB map[string]Favorite
	Favorite    struct {
		Name string `json:"name" binding:"required"`
	}

	TagsDB map[string]Tag
	Tag    mapset.Set

	PostsDB map[string]Post
	Post    struct {
		Tags      []string `json:"tags" binding:"required"`
		Namespace string   `json:"namespace" binding:"required"`
		External  bool     `json:"external"`
		ID        string   `json:"id" binding:"required"`
		Value     string   `json:"value" binding:"required"`
	}

	Data struct {
		*FavoritesDB
		*TagsDB
		*PostsDB
	}
)

var (
	favoritesDB = FavoritesDB{}
	tagsDB      = TagsDB{}
	postsDB     = PostsDB{}
)

func init() {
	router := gin.Default()
	http.Handle("/", router)

	v1Group := router.Group("/v1")

	v1Group.GET("/download", httpHandleDownloadDatabase)
	v1Group.POST("/upload", httpHandleUploadDatabase)

	v1Group.GET("/tags", httpHandleTagsIndex)

	favoritesGroup := v1Group.Group("/favorites")
	favoritesGroup.GET("/:key", httpHandleFavoriteRead)
	favoritesGroup.DELETE("/:key", httpHandleFavoriteDelete)
	favoritesGroup.GET("", httpHandleFavoriteIndex)
	favoritesGroup.POST("", httpHandleFavoriteCreate)

	postsGroup := v1Group.Group("/posts")
	postsGroup.GET("/:key", httpHandlePostRead)
	postsGroup.DELETE("/:key", httpHandlePostDelete)
	postsGroup.GET("", httpHandlePostIndex)
	postsGroup.POST("", httpHandlePostCreate)
}

func main() {
	appengine.Main()
	fmt.Println("")
}

func bindErrorResponse(err error) map[string][]string {
	errors := make(map[string][]string)
	switch msg := err.(type) {
	case validator.ValidationErrors:
		for _, v := range msg {
			errors[v.Field()] = append(errors[v.Field()], v.Tag())
		}
	}
	return errors

}
