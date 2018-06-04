package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpHandleFavoriteIndex(c *gin.Context) {
	favorites := databaseFavoritesQuery()
	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func httpHandleFavoriteCreate(c *gin.Context) {
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

func httpHandleFavoriteRead(c *gin.Context) {
	favorite, ok := databaseFavoriteRead(c.Param("key"))
	if !ok {
		c.String(http.StatusNotFound, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"favorite": favorite})
}

func httpHandleFavoriteDelete(c *gin.Context) {
	databaseFavoriteDelete(c.Param("key"))
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
