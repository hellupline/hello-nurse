package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpHandlePostIndex(c *gin.Context) {
	queryRaw := c.Query("q")

	var query interface{}
	json.Unmarshal([]byte(queryRaw), &query)

	posts := databasePostsQuery(query)
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func httpHandlePostCreate(c *gin.Context) {
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

func httpHandlePostRead(c *gin.Context) {
	post, ok := databasePostRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func httpHandlePostDelete(c *gin.Context) {
	databasePostDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
