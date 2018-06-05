package main

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpHandleUploadDatabase(c *gin.Context) {
	data := Data{}
	if err := gob.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}
	favoritesDB = *data.FavoritesDB
	tagsDB = *data.TagsDB
	postsDB = *data.PostsDB

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func httpHandleDownloadDatabase(c *gin.Context) {
	_ = gob.NewEncoder(c.Writer).Encode(Data{
		&favoritesDB,
		&tagsDB,
		&postsDB,
	})
}

func httpHandleTagsIndex(c *gin.Context) {
	tagNames := make([]string, 0, len(tagsDB))
	for key := range tagsDB {
		tagNames = append(tagNames, key)
	}
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}
