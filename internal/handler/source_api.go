package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sourceService interface {
	ListSources(ctx context.Context) ([]string, error)
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
