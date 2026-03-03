package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/Aquila-f/photo-slider/internal/domain"
	"github.com/Aquila-f/photo-slider/internal/photo"
	"github.com/gin-gonic/gin"
)

type albumService interface {
	ListAlbums(ctx context.Context) []domain.AlbumItem
	ListPhoto(ctx context.Context, albumKey string) ([]string, error)
	ReadPhoto(ctx context.Context, albumKey, photoToken string) ([]byte, error)
}

type AlbumAPI struct {
	svc          albumService
	compressor   photo.Compressor
	cacher       photo.Cacher
	extractor    domain.MetaExtractor
	listStrategy domain.PhotoListStrategy
}

func NewAlbumAPI(svc albumService, compressor photo.Compressor, cacher photo.Cacher, extractor domain.MetaExtractor, listStrategy domain.PhotoListStrategy) *AlbumAPI {
	return &AlbumAPI{svc: svc, compressor: compressor, cacher: cacher, extractor: extractor, listStrategy: listStrategy}
}

func (h *AlbumAPI) listAlbums(c *gin.Context) {
	albums := h.svc.ListAlbums(c.Request.Context())
	c.JSON(http.StatusOK, albums)
}

func (h *AlbumAPI) listPhotos(c *gin.Context) {
	albumKey := c.Param("albumkey")
	photos, err := h.svc.ListPhoto(c.Request.Context(), albumKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if c.Query("shuffle") == "true" {
		photos, _ = h.listStrategy.Arrange(c.Request.Context(), photos)
	}
	c.JSON(http.StatusOK, photos)
}

func (h *AlbumAPI) readPhoto(c *gin.Context) {
	albumKey := c.Param("albumkey")
	token := c.Param("key")
	ctx := c.Request.Context()

	cacheKey := albumKey + "/" + token
	if cached, err := h.cacher.Get(ctx, cacheKey); err == nil {
		log.Printf("cache hit: %s", cacheKey)
		setMetaHeaders(c, cached.Meta)
		c.Data(http.StatusOK, http.DetectContentType(cached.Data), cached.Data)
		return
	}

	raw, err := h.svc.ReadPhoto(ctx, albumKey, token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photo not found"})
		return
	}

	// meta is best-effort; failure just means no EXIF headers in the response.
	meta, _ := h.extractor.Extract(ctx, raw)

	data, err := h.compressor.Compress(ctx, raw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = h.cacher.Set(ctx, cacheKey, photo.CachedPhoto{Data: data, Meta: meta})

	setMetaHeaders(c, meta)
	c.Data(http.StatusOK, http.DetectContentType(data), data)
}

func setMetaHeaders(c *gin.Context, meta *domain.PhotoMeta) {
	if meta == nil {
		return
	}
	for k, v := range meta.Headers() {
		c.Header(k, v)
	}
}
