package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(indexHTML string, api *AlbumAPI) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})

	r.GET("/api/albums", api.listAlbums)
	r.GET("/api/albums/:album", api.listPhotos)
	r.GET("/photos/:album/:key", api.readPhoto)

	return r
}
