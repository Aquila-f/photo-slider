package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sourceRequest struct {
	ID string `json:"id" binding:"required"`
}

type sourceService interface {
	ListSources(ctx context.Context) ([]string, error)
	AddSource(ctx context.Context, id string) error
	DeleteSource(ctx context.Context, id string) error
}

type SourceAPI struct {
	svc sourceService
}

func NewSourceAPI(svc sourceService) *SourceAPI {
	return &SourceAPI{svc: svc}
}

func (h *SourceAPI) listSources(c *gin.Context) {
	sources, err := h.svc.ListSources(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sources)
}

func (h *SourceAPI) createSource(c *gin.Context) {
	var req sourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AddSource(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *SourceAPI) deleteSource(c *gin.Context) {
	var req sourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.DeleteSource(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
