package main

import (
	"net/http"
	"sort"

	"github.com/hellupline/hello-nurse/nursetags"

	"github.com/gin-gonic/gin"
)

func HttpHandlePostIndex(c *gin.Context) {
	posts := database.PostIndex(c.Query("q"))

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Key < posts[j].Key
	})
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func HttpHandlePostCreate(c *gin.Context) {
	post := nursetags.PostData{}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	database.PostCreate(post)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandlePostRead(c *gin.Context) {
	key := nursetags.PostKey{
		Namespace: c.Param("namespace"),
		Key:       c.Param("key"),
	}
	post, ok := database.PostRead(key)
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func HttpHandlePostDelete(c *gin.Context) {
	key := nursetags.PostKey{
		Namespace: c.Param("namespace"),
		Key:       c.Param("key"),
	}
	database.PostDelete(key)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleTagsIndex(c *gin.Context) {
	tagNames := database.TagIndex()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i].Count > tagNames[j].Count
	})
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}
