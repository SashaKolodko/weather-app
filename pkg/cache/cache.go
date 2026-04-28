package cache

import (
    "time"
)

// CacheItem представляет элемент кэша
type CacheItem struct {
    Data      []byte    `json:"data"`
    Timestamp time.Time `json:"timestamp"`
    ExpiresAt time.Time `json:"expires_at"`
}

// Cache интерфейс для кэширования
type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, data []byte, duration time.Duration)
    Delete(key string)
    Clear()
}