package cache

import (
    "encoding/json"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// FileCache реализует файловое кэширование
type FileCache struct {
    cacheDir string
    mu       sync.RWMutex
}

// NewFileCache создает новый файловый кэш
func NewFileCache(cacheDir string) (*FileCache, error) {
    if cacheDir == "" {
        cacheDir = filepath.Join(os.TempDir(), "weather-cache")
    }
    
    if err := os.MkdirAll(cacheDir, 0755); err != nil {
        return nil, err
    }
    
    return &FileCache{
        cacheDir: cacheDir,
    }, nil
}

// Get получает данные из файлового кэша
func (c *FileCache) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    filePath := filepath.Join(c.cacheDir, key+".json")
    
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, false
    }
    
    var item CacheItem
    if err := json.Unmarshal(data, &item); err != nil {
        return nil, false
    }
    
    // Проверяем срок действия
    if time.Now().After(item.ExpiresAt) {
        os.Remove(filePath)
        return nil, false
    }
    
    return item.Data, true
}

// Set сохраняет данные в файловый кэш
func (c *FileCache) Set(key string, data []byte, duration time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    item := CacheItem{
        Data:      data,
        Timestamp: time.Now(),
        ExpiresAt: time.Now().Add(duration),
    }
    
    fileData, err := json.Marshal(item)
    if err != nil {
        return
    }
    
    filePath := filepath.Join(c.cacheDir, key+".json")
    os.WriteFile(filePath, fileData, 0644)
}

// Delete удаляет файл из кэша
func (c *FileCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    filePath := filepath.Join(c.cacheDir, key+".json")
    os.Remove(filePath)
}

// Clear очищает всю директорию кэша
func (c *FileCache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    files, _ := filepath.Glob(filepath.Join(c.cacheDir, "*.json"))
    for _, file := range files {
        os.Remove(file)
    }
}