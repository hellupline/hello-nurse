package main

import (
	"net/http"
	"sort"

	"github.com/hellupline/hello-nurse/nursetags"

	"github.com/gin-gonic/gin"
)

func HttpHandleFavoriteIndex(c *gin.Context) {
	favorites := nursetags.DefaultDatabase.FavoriteQuery()
	sort.Slice(favorites, func(i, j int) bool {
		return favorites[i].Name < favorites[j].Name
	})
	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func HttpHandleFavoriteCreate(c *gin.Context) {
	favorite := nursetags.Favorite{}
	if err := c.ShouldBindJSON(&favorite); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": bindErrorResponse(err),
			"result": "error",
		})
		return
	}

	nursetags.DefaultDatabase.FavoriteCreate(favorite)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleFavoriteRead(c *gin.Context) {
	favorite, ok := nursetags.DefaultDatabase.FavoriteRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"favorite": favorite})
}

func HttpHandleFavoriteDelete(c *gin.Context) {
	nursetags.DefaultDatabase.FavoriteDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleTagsIndex(c *gin.Context) {
	tagNames := nursetags.DefaultDatabase.TagQuery()
	sort.Slice(tagNames, func(i, j int) bool {
		return tagNames[i].Count > tagNames[j].Count
	})
	c.JSON(http.StatusOK, gin.H{"tags": tagNames})
}

func HttpHandlePostIndex(c *gin.Context) {
	posts := nursetags.DefaultDatabase.PostQuery(c.Query("q"))

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

	nursetags.DefaultDatabase.PostCreate(post)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandlePostRead(c *gin.Context) {
	key := nursetags.PostKey{
		Namespace: c.Param("namespace"),
		Key:       c.Param("key"),
	}
	post, ok := nursetags.DefaultDatabase.PostRead(key)
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
	nursetags.DefaultDatabase.PostDelete(key)
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDatasetUploadGOB(c *gin.Context) {
	if err := nursetags.DefaultDatabase.Read(c.Request.Body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func HttpHandleDatasetDownloadGOB(c *gin.Context) {
	_ = nursetags.DefaultDatabase.Write(c.Writer)
}
