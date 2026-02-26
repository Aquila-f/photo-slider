package photo

import (
	"context"
	"errors"
	"sync"
)

var ErrCacheMiss = errors.New("cache miss")

type Cacher interface {
	Set(ctx context.Context, token string, data []byte) error
	Get(ctx context.Context, token string) ([]byte, error)
}


type FixedSizeMapCacher struct {
	mu       sync.RWMutex
	size     int
	idx      int
	indexMap []string
	store    map[string][]byte
}

func NewFixedSizeMapCacher(size int) *FixedSizeMapCacher {
	if size <= 0 {
		panic("FixedSizeMapCacher: size must be greater than 0")
	}
	return &FixedSizeMapCacher{
		size:     size,
		indexMap: make([]string, size),
		store:    make(map[string][]byte, size),
	}
}

func (c *FixedSizeMapCacher) Set(_ context.Context, token string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exist := c.store[token]; exist {
		c.store[token] = data
		return nil
	}

	if old := c.indexMap[c.idx]; old != "" {
		delete(c.store, old)
	}
	c.indexMap[c.idx] = token
	c.store[token] = data
	c.idx = (c.idx + 1) % c.size
	return nil
}

func (c *FixedSizeMapCacher) Get(_ context.Context, token string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.store[token]
	if !ok {
		return nil, ErrCacheMiss
	}
	return data, nil
}
	
