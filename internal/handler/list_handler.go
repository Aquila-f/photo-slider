package handler

import (
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	source photo.Source
}

func NewListHandler(s photo.Source) *ListHandler {
	return &ListHandler{source: s}
}

func (h *ListHandler) Handle(c *gin.Context) {
	images, err := h.source.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, images)
}
