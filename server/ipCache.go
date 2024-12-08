package server

import (
	"sync"

	"github.com/google/uuid"
)

type IPCache struct {
	ipToClientID map[string]uuid.UUID
	mu           sync.Mutex
}

func NewIPCache() *IPCache {
	return &IPCache{
		ipToClientID: make(map[string]uuid.UUID),
	}
}

func (c *IPCache) GetByIP(ip string) (uuid.UUID, bool) {
	id, found := c.ipToClientID[ip]
	return id, found
}

func (c *IPCache) AddIP(ip string, id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ipToClientID[ip] = id
}
