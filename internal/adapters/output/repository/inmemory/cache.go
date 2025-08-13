package inmemory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/shanth1/authorization/internal/core/domain"
)

type cacheItem struct {
	Value      []byte
	Expiration int64
}

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]cacheItem),
	}
}

func (c *InMemoryCache) SetState(ctx context.Context, key string, expiration time.Duration) error {
	return c.set(key, []byte("pending"), expiration)
}

func (c *InMemoryCache) CheckState(ctx context.Context, key string) error {
	_, err := c.get(key)
	return err
}

func (c *InMemoryCache) DeleteState(ctx context.Context, key string) error {
	c.del(key)
	return nil
}

func (c *InMemoryCache) SetSession(ctx context.Context, key string, tokens *domain.Tokens, expiration time.Duration) error {
	data, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	return c.set(key, data, expiration)
}

func (c *InMemoryCache) GetSession(ctx context.Context, key string) (*domain.Tokens, error) {
	data, err := c.get(key)
	if err != nil {
		return nil, err
	}
	var tokens domain.Tokens
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	return &tokens, nil
}

func (c *InMemoryCache) DeleteSession(ctx context.Context, key string) error {
	c.del(key)
	return nil
}

func (c *InMemoryCache) set(key string, value []byte, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration).UnixNano(),
	}
	return nil
}

func (c *InMemoryCache) get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, fmt.Errorf("key not found in cache")
	}
	if time.Now().UnixNano() > item.Expiration {
		return nil, fmt.Errorf("key expired")
	}
	return item.Value, nil
}

func (c *InMemoryCache) del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}
