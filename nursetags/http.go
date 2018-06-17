package nursetags

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator"
)

type DatabaseDump struct {
	Favorites FavoritesDB `json:"favorites"`
	Tags      TagsDB      `json:"tags"`
	Posts     PostsDB     `json:"posts"`
}

func HttpHandleUploadDatabaseJSON(c *gin.Context) {
	var data DatabaseDump
	database.Lock()
	defer database.Unlock()

	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	database.Favorites = data.Favorites
	database.Tags = data.Tags
	database.Posts = data.Posts

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDownloadDatabaseJSON(c *gin.Context) {
	database.RLock()
	defer database.RUnlock()

	_ = json.NewEncoder(c.Writer).Encode(DatabaseDump{
		Favorites: database.Favorites,
		Tags:      database.Tags,
		Posts:     database.Posts,
	})
}

func HttpHandleUploadDatabaseGOB(c *gin.Context) {
	var data DatabaseDump
	database.Lock()
	defer database.Unlock()

	if err := gob.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	database.Favorites = data.Favorites
	database.Tags = data.Tags
	database.Posts = data.Posts

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDownloadDatabaseGOB(c *gin.Context) {
	database.RLock()
	defer database.RUnlock()

	_ = gob.NewEncoder(c.Writer).Encode(DatabaseDump{
		Favorites: database.Favorites,
		Tags:      database.Tags,
		Posts:     database.Posts,
	})
}

func HttpHandleFavoriteIndex(c *gin.Context) {
	favorites := databaseFavoritesQuery()
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

	databaseFavoriteCreate(favorite)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleFavoriteRead(c *gin.Context) {
	favorite, ok := databaseFavoriteRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"favorite": favorite})
}

func HttpHandleFavoriteDelete(c *gin.Context) {
	databaseFavoriteDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleTagsIndex(c *gin.Context) {
	tagNames := databaseTagsQuery()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i] < tagNames[j]
	})
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}

func HttpHandlePostIndex(c *gin.Context) {
	posts := databasePostsQuery(c.Query("q"))

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

	databasePostCreate(post)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandlePostRead(c *gin.Context) {
	key := PostKey{c.Param("namespace"), c.Param("key")}
	post, ok := databasePostRead(key)
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func HttpHandlePostDelete(c *gin.Context) {
	key := PostKey{c.Param("namespace"), c.Param("key")}
	databasePostDelete(key)
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