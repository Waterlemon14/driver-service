package adapter

import (
	"context"
	"driver-service/internal/domain"
	"fmt"
	"sync"
	"time"
)

// --- In-Memory Cache ---
type MemCache struct {
	store sync.Map // Key: string, Value: cacheItem
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewMemCache() *MemCache {
	return &MemCache{}
}

func (c *MemCache) GetList(ctx context.Context, page, limit int) (*domain.ListResponse, bool) {
	key := fmt.Sprintf("list:%d:%d", page, limit)
	val, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}
	item := val.(cacheItem)
	if time.Now().After(item.expiresAt) {
		c.store.Delete(key)
		return nil, false
	}
	return item.value.(*domain.ListResponse), true
}

func (c *MemCache) SetList(ctx context.Context, page, limit int, data *domain.ListResponse) {
	key := fmt.Sprintf("list:%d:%d", page, limit)
	c.store.Store(key, cacheItem{
		value:     data,
		expiresAt: time.Now().Add(30 * time.Second), // TTL 30s
	})
}

func (c *MemCache) InvalidateList(ctx context.Context) {
	// clear everything
	c.store.Range(func(key, value interface{}) bool {
		c.store.Delete(key)
		return true
	})
}
