package handler

import (
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/gin-gonic/gin"
)

type ReadHandler struct {
	resolver photo.Resolver
	reader   photo.Reader
}

func NewReadHandler(r photo.Resolver, rd photo.Reader) *ReadHandler {
	return &ReadHandler{resolver: r, reader: rd}
}

func (h *ReadHandler) Handle(c *gin.Context) {
	token := c.Param("key")
	path, err := h.resolver.Resolve(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	data, err := h.reader.Read(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	contentType := http.DetectContentType(data)
	c.Data(http.StatusOK, contentType, data)
}
