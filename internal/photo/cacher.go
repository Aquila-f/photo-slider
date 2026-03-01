package photo

import (
	"context"
	"errors"
	"sync"
)

var ErrCacheMiss = errors.New("cache miss")

type Cacher interface {
	Set(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type FixedSizeMapCacher struct {
	mu    sync.RWMutex
	size  int
	head  int
	ring  []string
	store map[string][]byte
}

func NewFixedSizeMapCacher(size int) *FixedSizeMapCacher {
	if size <= 0 {
		panic("FixedSizeMapCacher: size must be greater than 0")
	}
	return &FixedSizeMapCacher{
		size:  size,
		ring:  make([]string, size),
		store: make(map[string][]byte, size),
	}
}

func (c *FixedSizeMapCacher) Set(_ context.Context, key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.store[key]; ok {
		c.store[key] = data
		return nil
	}

	if old := c.ring[c.head]; old != "" {
		delete(c.store, old)
	}
	c.ring[c.head] = key
	c.store[key] = data
	c.head = (c.head + 1) % c.size
	return nil
}

func (c *FixedSizeMapCacher) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.store[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	return data, nil
}
