package redis

import (
	"github.com/go-redis/cache/v9"
)

type Cache struct {
	rc *cache.Cache
}

func NewCache(rc *cache.Cache) *Cache {
	return &Cache{
		rc: rc,
	}
}

func (c *Cache) Once(item *cache.Item) error {
	return c.rc.Once(item)
}
