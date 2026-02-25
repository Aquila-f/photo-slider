package handler

import (
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/lister"
	"github.com/Aquila-f/photo-slider/internal/reader"
	"github.com/gin-gonic/gin"
)

type ReadHandler struct {
	lister lister.FileLister
	reader reader.PhotoReader
}

func NewReadHandler(l lister.FileLister, r reader.PhotoReader) *ReadHandler {
	return &ReadHandler{lister: l, reader: r}
}

func (h *ReadHandler) Handle(c *gin.Context) {
	key := c.Param("key")
	path, err := h.lister.GetCompletePath(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	data, err := h.reader.ReadPhoto(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	contentType := http.DetectContentType(data)
	c.Data(http.StatusOK, contentType, data)
}
