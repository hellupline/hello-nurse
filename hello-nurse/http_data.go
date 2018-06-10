package main

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpHandleUploadDatabase(c *gin.Context) {
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

func httpHandleDownloadDatabase(c *gin.Context) {
	database.RLock()
	_ = gob.NewEncoder(c.Writer).Encode(database)
	database.RUnlock()
}

func httpHandleTagsIndex(c *gin.Context) {
	tagNames := databaseTagsQuery()
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}
