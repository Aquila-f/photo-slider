package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(staticFS embed.FS, api *AlbumAPI, sourceAPI *SourceAPI) *gin.Engine {
	r := gin.Default()

	// Serve index.html at root
	indexHTML, _ := staticFS.ReadFile("static/index.html")
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	// Serve static assets (CSS, JS)
	staticSub, _ := fs.Sub(staticFS, "static")
	r.StaticFS("/static", http.FS(staticSub))

	r.GET("/api/sources", sourceAPI.listSources)
	r.GET("/api/albums", api.listAlbums)
	r.GET("/api/albums/:albumkey", api.listPhotos)
	r.GET("/photos/:albumkey/:key", api.readPhoto)

	return r
}
