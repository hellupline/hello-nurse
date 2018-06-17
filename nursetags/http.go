package nursetags

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator"
)

type DatabaseDump struct {
	Favorites FavoritesDB
	Tags      TagsDB
	Posts     PostsDB
}

func HttpHandleUploadDatabaseGOB(c *gin.Context) {
	if err := database.Read(c.Request.Body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDownloadDatabaseGOB(c *gin.Context) {
	_ = database.Write(c.Writer)
}

func HttpHandleFavoriteIndex(c *gin.Context) {
	favorites := DatabaseFavoritesQuery()
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

	DatabaseFavoriteCreate(favorite)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleFavoriteRead(c *gin.Context) {
	favorite, ok := DatabaseFavoriteRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"favorite": favorite})
}

func HttpHandleFavoriteDelete(c *gin.Context) {
	DatabaseFavoriteDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleTagsIndex(c *gin.Context) {
	tagNames := DatabaseTagsQuery()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i] < tagNames[j]
	})
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}

func HttpHandlePostIndex(c *gin.Context) {
	posts := DatabasePostsQuery(c.Query("q"))

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID < posts[j].ID
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

	DatabasePostCreate(post)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandlePostRead(c *gin.Context) {
	key := PostKey{c.Param("namespace"), c.Param("key")}
	post, ok := DatabasePostRead(key)
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func HttpHandlePostDelete(c *gin.Context) {
	key := PostKey{c.Param("namespace"), c.Param("key")}
	DatabasePostDelete(key)
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
