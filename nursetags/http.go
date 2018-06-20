package nursetags

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator"
)

func HttpHandleUploadDatasetGOB(c *gin.Context) {
	if err := DefaultDatabase.Read(c.Request.Body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDownloadDatasetGOB(c *gin.Context) {
	_ = DefaultDatabase.Write(c.Writer)
}

func HttpHandleFavoriteIndex(c *gin.Context) {
	favorites := DefaultDatabase.FavoriteQuery()
	sort.Slice(favorites, func(i, j int) bool {
		return favorites[i].Name < favorites[j].Name
	})
	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func HttpHandleFavoriteCreate(c *gin.Context) {
	favorite := Favorite{}
	if err := c.ShouldBindJSON(&favorite); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	DefaultDatabase.FavoriteCreate(favorite)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleFavoriteRead(c *gin.Context) {
	favorite, ok := DefaultDatabase.FavoriteRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"favorite": favorite})
}

func HttpHandleFavoriteDelete(c *gin.Context) {
	DefaultDatabase.FavoriteDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleTagsIndex(c *gin.Context) {
	tagNames := DefaultDatabase.TagQuery()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i] < tagNames[j]
	})
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}

func HttpHandlePostIndex(c *gin.Context) {
	posts := DefaultDatabase.PostQuery(c.Query("q"))

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Key < posts[j].Key
	})
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func HttpHandlePostCreate(c *gin.Context) {
	post := PostData{}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	DefaultDatabase.PostCreate(post)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandlePostRead(c *gin.Context) {
	key := PostKey{c.Param("namespace"), c.Param("key")}
	post, ok := DefaultDatabase.PostRead(key)
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func HttpHandlePostDelete(c *gin.Context) {
	key := PostKey{c.Param("namespace"), c.Param("key")}
	DefaultDatabase.PostDelete(key)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
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
