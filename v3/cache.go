package v3

import (
	"errors"
	"sync"

	prismgoclient "github.com/nutanix-cloud-native/prism-go-client"
)

type clientCacheMap map[string]*Client

// ErrorClientNotFound is returned when the client is not found in the cache
var ErrorClientNotFound = errors.New("client not found in client cache")

// ClientCache is a cache for prism clients
type ClientCache struct {
	cache clientCacheMap
	mtx   sync.RWMutex
}

// NewClientCache returns a new ClientCache
func NewClientCache() *ClientCache {
	return &ClientCache{
		cache: make(clientCacheMap),
		mtx:   sync.RWMutex{},
	}
}

// NewCachedClient creates a new client and adds it to the cache
func (c *ClientCache) NewCachedClient(clientName string, credentials prismgoclient.Credentials, opts ...ClientOption) (*Client, error) {
	client, err := NewV3Client(credentials, opts...)
	if err != nil {
		return nil, err
	}

	c.Set(clientName, client)
	return client, nil
}

// Get returns the client for the given client name
func (c *ClientCache) Get(clientName string) (*Client, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	clnt, ok := c.cache[clientName]
	if !ok {
		return nil, ErrorClientNotFound
	}

	return clnt, nil
}

// Set adds the client to the cache
func (c *ClientCache) Set(clientName string, client *Client) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cache[clientName] = client
}

// Delete removes the client from the cache
func (c *ClientCache) Delete(clientName string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.cache, clientName)
}
