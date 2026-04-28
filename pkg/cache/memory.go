package cache

import (
    "sync"
    "time"
)

// MemoryCache реализует in-memory кэш
type MemoryCache struct {
    data map[string]CacheItem
    mu   sync.RWMutex
}

// NewMemoryCache создает новый in-memory кэш
func NewMemoryCache() *MemoryCache {
    cache := &MemoryCache{
        data: make(map[string]CacheItem),
    }
    
    // Запускаем горутину для очистки устаревших записей
    go cache.cleanup()
    
    return cache
}

// Get получает данные из кэша
func (c *MemoryCache) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.data[key]
    if !exists {
        return nil, false
    }
    
    // Проверяем не истек ли срок действия
    if time.Now().After(item.ExpiresAt) {
        return nil, false
    }
    
    return item.Data, true
}

// Set сохраняет данные в кэш
func (c *MemoryCache) Set(key string, data []byte, duration time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data[key] = CacheItem{
        Data:      data,
        Timestamp: time.Now(),
        ExpiresAt: time.Now().Add(duration),
    }
}

// Delete удаляет запись из кэша
func (c *MemoryCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    delete(c.data, key)
}

// Clear очищает весь кэш
func (c *MemoryCache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data = make(map[string]CacheItem)
}

// cleanup периодически очищает устаревшие записи
func (c *MemoryCache) cleanup() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.data {
            if now.After(item.ExpiresAt) {
                delete(c.data, key)
            }
        }
        c.mu.Unlock()
    }
}