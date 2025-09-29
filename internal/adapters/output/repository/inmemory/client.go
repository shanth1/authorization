package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/shanth1/authorization/internal/core/domain"
	"github.com/shanth1/gotools/errs"
)

type clientRepo struct {
	mu      sync.RWMutex
	clients map[domain.ClientID]*domain.Client
}

func NewClientRepo() *clientRepo {
	return &clientRepo{
		clients: make(map[domain.ClientID]*domain.Client),
	}
}

func (r *clientRepo) Get(ctx context.Context, id domain.ClientID) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[id]
	if !exists {
		return nil, errs.ErrNotFound
	}

	clientCopy := *client
	return &clientCopy, nil
}

func (r *clientRepo) Save(ctx context.Context, c *domain.Client) error {
	if c.ID == "" {
		return errors.New("emtpy id")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	clientToSave := *c
	r.clients[c.ID] = &clientToSave

	return nil
}
