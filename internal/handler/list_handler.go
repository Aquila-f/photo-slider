package handler

import (
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/lister"
	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	lister lister.FileLister
}

func NewListHandler(l lister.FileLister) *ListHandler {
	return &ListHandler{lister: l}
}

func (h *ListHandler) Handle(c *gin.Context) {
	images, err := h.lister.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, images)
}
