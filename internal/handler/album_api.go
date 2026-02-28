package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/gin-gonic/gin"
)

type albumService interface {
	ListAlbums(ctx context.Context) []string
	ListPhoto(ctx context.Context, albumName string) ([]string, error)
	ReadPhoto(ctx context.Context, albumName, photoToken string) ([]byte, error)
}

type AlbumAPI struct {
	svc        albumService
	compressor photo.Compressor
	cacher     photo.Cacher
}

func NewAlbumAPI(svc albumService, compressor photo.Compressor, cacher photo.Cacher) *AlbumAPI {
	return &AlbumAPI{svc: svc, compressor: compressor, cacher: cacher}
}

func (h *AlbumAPI) listAlbums(c *gin.Context) {
	albums := h.svc.ListAlbums(c.Request.Context())
	c.JSON(http.StatusOK, albums)
}

func (h *AlbumAPI) listPhotos(c *gin.Context) {
	albumName := c.Param("album")
	photos, err := h.svc.ListPhoto(c.Request.Context(), albumName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, photos)
}

func (h *AlbumAPI) readPhoto(c *gin.Context) {
	albumName := c.Param("album")
	token := c.Param("key")
	ctx := c.Request.Context()

	if data, err := h.cacher.Get(ctx, token); err == nil {
		log.Printf("cache hit: %s", token)
		c.Data(http.StatusOK, http.DetectContentType(data), data)
		return
	}

	data, err := h.svc.ReadPhoto(ctx, albumName, token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}
	data, err = h.compressor.Compress(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = h.cacher.Set(ctx, token, data)
	c.Data(http.StatusOK, http.DetectContentType(data), data)
}
