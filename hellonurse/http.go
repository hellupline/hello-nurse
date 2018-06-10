package main

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func HttpHandleUploadDatabase(c *gin.Context) {
	var data Database
	if err := gob.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	database.RLock()
	database.Favorites = data.Favorites
	database.Tags = data.Tags
	database.Posts = data.Posts
	database.RUnlock()

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDownloadDatabase(c *gin.Context) {
	database.RLock()
	_ = gob.NewEncoder(c.Writer).Encode(database)
	database.RUnlock()
}

func HttpHandleFavoriteIndex(c *gin.Context) {
	favorites := databaseFavoritesQuery()
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
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}

func HttpHandlePostIndex(c *gin.Context) {
	queryRaw := c.Query("q")

	var query interface{}
	json.Unmarshal([]byte(queryRaw), &query)

	posts := databasePostsQuery(query)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID < posts[j].ID
	})
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func HttpHandlePostCreate(c *gin.Context) {
	post := Post{}
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
	post, ok := databasePostRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func HttpHandlePostDelete(c *gin.Context) {
	databasePostDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
