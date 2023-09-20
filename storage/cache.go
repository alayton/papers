package storage

import (
	"time"

	"github.com/erni27/imcache"

	"github.com/alayton/papers"
)

func NewTTLCache(cleanupInterval time.Duration) *TTLCache {
	return &TTLCache{
		cache: imcache.New[string, string](imcache.WithCleanerOption[string, string](cleanupInterval)),
	}
}

type TTLCache struct {
	cache *imcache.Cache[string, string]
}

func (c TTLCache) Set(key, value string, ttl time.Duration) {
	c.cache.Set(key, value, imcache.WithExpiration(ttl))
}

func (c TTLCache) Get(key string) (string, bool) {
	return c.cache.Get(key)
}

func NewTokenCache(p *papers.Papers) *TokenCache {
	return &TokenCache{
		papers: p,
		cache:  imcache.New[string, bool](imcache.WithCleanerOption[string, bool](p.Config.TokenCacheExpiration * 2)),
	}
}

type TokenCache struct {
	papers *papers.Papers
	cache  *imcache.Cache[string, bool]
}

// Checks if a token is valid. Chains should only be cached when they're invalidated, which takes precedence over an individual token
func (c TokenCache) IsTokenValid(token, chain string) (bool, bool) {
	chainValid, chainExists := c.cache.Get(chain)
	if chainExists && !chainValid {
		return false, true
	}
	return c.cache.Get(token)
}

func (c TokenCache) SetTokenValidity(token string, valid bool) {
	c.cache.Set(token, valid, imcache.WithExpiration(c.papers.Config.TokenCacheExpiration))
}

func (c TokenCache) SetChainValidity(chain string, valid bool) {
	c.cache.Set(chain, valid, imcache.WithExpiration(c.papers.Config.TokenCacheExpiration))
}
