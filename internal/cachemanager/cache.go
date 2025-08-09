package cachemanager

import (
	"context"
	"log"

	"github.com/fengziwei0826/caching_proxy_exe/pkg/db"
)

const (
	cacheKey = "cache_key"
)

type CacheManager interface {
	CacheResponse(key, value string) error
	GetResponse(key string) (string, error)
	ClearCache(keys ...string) error
	ClearAllCache() error
	Close() error
}

type cacheManager struct {
	ctx   context.Context
	cache db.CacheProxy
}

func (c *cacheManager) CacheResponse(key, value string) error {
	err := c.cache.SAdd(cacheKey, key)
	if err != nil {
		log.Printf("Failed to add key to cache set: %v", err)
		return err
	}
	return c.cache.Set(key, value)
}

func (c *cacheManager) GetResponse(key string) (string, error) {
	return c.cache.Get(key)
}

func (c *cacheManager) ClearCache(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	_, err := c.cache.Del(keys...)
	return err
}

func (c *cacheManager) ClearAllCache() error {
	keys, err := c.cache.GetSet(cacheKey)
	if err != nil {
		log.Printf("Failed to get cache keys: %v", err)
		return err
	}
	if len(keys) == 0 {
		log.Println("No cache keys to clear.")
		return err
	}
	n, err := c.cache.Del(keys...)
	if err != nil {
		log.Printf("Failed to delete cache keys: %v", err)
		return err
	}
	c.cache.Del(cacheKey)
	log.Printf("Cleared %d cache keys.", n)
	return nil
}

func (c *cacheManager) Close() error {
	return c.cache.Close()
}

func NewCacheManager(ctx context.Context, cache db.CacheProxy) CacheManager {
	return &cacheManager{
		ctx:   ctx,
		cache: cache,
	}
}
