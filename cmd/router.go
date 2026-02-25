package main

import (
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/handler"
	"github.com/gin-gonic/gin"
)

func setupRouter(lh *handler.ListHandler, rh *handler.ReadHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})

	r.GET("/api/images", lh.Handle)
	r.GET("/images/:key", rh.Handle)

	return r
}
