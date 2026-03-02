package photo

import (
	"context"
	"errors"
	"sync"

	"github.com/Aquila-f/photo-slider/internal/domain"
)

var ErrCacheMiss = errors.New("cache miss")

type CachedPhoto struct {
	Data []byte
	Meta *domain.PhotoMeta
}

type Cacher interface {
	Set(ctx context.Context, key string, photo CachedPhoto) error
	Get(ctx context.Context, key string) (CachedPhoto, error)
}

type FixedSizeMapCacher struct {
	mu    sync.RWMutex
	size  int
	head  int
	ring  []string
	store map[string]CachedPhoto
}

func NewFixedSizeMapCacher(size int) *FixedSizeMapCacher {
	if size <= 0 {
		panic("FixedSizeMapCacher: size must be greater than 0")
	}
	return &FixedSizeMapCacher{
		size:  size,
		ring:  make([]string, size),
		store: make(map[string]CachedPhoto, size),
	}
}

func (c *FixedSizeMapCacher) Set(_ context.Context, key string, photo CachedPhoto) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.store[key]; ok {
		c.store[key] = photo
		return nil
	}

	if old := c.ring[c.head]; old != "" {
		delete(c.store, old)
	}
	c.ring[c.head] = key
	c.store[key] = photo
	c.head = (c.head + 1) % c.size
	return nil
}

func (c *FixedSizeMapCacher) Get(_ context.Context, key string) (CachedPhoto, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	photo, ok := c.store[key]
	if !ok {
		return CachedPhoto{}, ErrCacheMiss
	}
	return photo, nil
}
