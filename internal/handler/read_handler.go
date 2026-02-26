package handler

import (
	"log"
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/gin-gonic/gin"
)

type ReadHandler struct {
	resolver   photo.Resolver
	reader     photo.Reader
	compressor photo.Compressor
	cacher     photo.Cacher
}

func NewReadHandler(r photo.Resolver, rd photo.Reader, c photo.Compressor, ca photo.Cacher) *ReadHandler {
	return &ReadHandler{resolver: r, reader: rd, compressor: c, cacher: ca}
}

func (h *ReadHandler) Handle(c *gin.Context) {
	token := c.Param("key")
	ctx := c.Request.Context()

	if data, err := h.cacher.Get(ctx, token); err == nil {
		log.Printf("cache hit: %s", token)
		contentType := http.DetectContentType(data)
		c.Data(http.StatusOK, contentType, data)
		return
	}

	path, err := h.resolver.Resolve(ctx, token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	data, err := h.reader.Read(ctx, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data, err = h.compressor.Compress(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = h.cacher.Set(ctx, token, data)
	contentType := http.DetectContentType(data)
	c.Data(http.StatusOK, contentType, data)
}
