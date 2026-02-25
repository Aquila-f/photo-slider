package handler

import (
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	resolver photo.Resolver
}

func NewListHandler(r photo.Resolver) *ListHandler {
	return &ListHandler{resolver: r}
}

func (h *ListHandler) Handle(c *gin.Context) {
	images, err := h.resolver.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, images)
}
