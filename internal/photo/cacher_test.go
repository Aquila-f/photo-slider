package photo

import (
	"context"
	"errors"
	"sync"
	"testing"
)

var ctx = context.Background()

func TestFixedSizeMapCacher_GetMiss(t *testing.T) {
	c := NewFixedSizeMapCacher(2)
	_, err := c.Get(ctx, "missing")
	if !errors.Is(err, ErrCacheMiss) {
		t.Errorf("Get() error = %v, want ErrCacheMiss", err)
	}
}

func TestFixedSizeMapCacher_SetAndGet(t *testing.T) {
	c := NewFixedSizeMapCacher(2)
	data := []byte("hello")

	if err := c.Set(ctx, "a", data); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	got, err := c.Get(ctx, "a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("Get() = %q, want %q", got, data)
	}
}

func TestFixedSizeMapCacher_UpdateExisting(t *testing.T) {
	c := NewFixedSizeMapCacher(2)
	_ = c.Set(ctx, "a", []byte("old"))
	_ = c.Set(ctx, "a", []byte("new"))

	got, err := c.Get(ctx, "a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(got) != "new" {
		t.Errorf("Get() = %q, want %q", got, "new")
	}
}

func TestFixedSizeMapCacher_Eviction(t *testing.T) {
	c := NewFixedSizeMapCacher(2)
	_ = c.Set(ctx, "a", []byte("a"))
	_ = c.Set(ctx, "b", []byte("b"))
	_ = c.Set(ctx, "c", []byte("c"))

	if _, err := c.Get(ctx, "a"); !errors.Is(err, ErrCacheMiss) {
		t.Error("Get(a) expected ErrCacheMiss after eviction")
	}
	if _, err := c.Get(ctx, "b"); err != nil {
		t.Errorf("Get(b) error = %v, want hit", err)
	}
	if _, err := c.Get(ctx, "c"); err != nil {
		t.Errorf("Get(c) error = %v, want hit", err)
	}
}

func TestFixedSizeMapCacher_InvalidSize(t *testing.T) {
	for _, size := range []int{0, -1} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("NewFixedSizeMapCacher(%d) expected panic, got none", size)
				}
			}()
			NewFixedSizeMapCacher(size)
		}()
	}
}

func TestFixedSizeMapCacher_ConcurrentAccess(t *testing.T) {
	c := NewFixedSizeMapCacher(10)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			token := string(rune('a' + i%10))
			_ = c.Set(ctx, token, []byte(token))
			_, _ = c.Get(ctx, token)
		}(i)
	}
	wg.Wait()
}
