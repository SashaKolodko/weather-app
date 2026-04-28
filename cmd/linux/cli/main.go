package main

import (
    "os"
    
    "weather-app/internal/adapters/weather"
    "weather-app/internal/pkg/app/cli"
    "weather-app/pkg/cache"
    "weather-app/pkg/logger"
)

func main() {
    log := logger.New()
    
    weatherCache, err := cache.NewFileCache("")
    if err != nil {
        log.Error("Failed to create cache", err)
        os.Exit(1)
    }
    
    log.Info("Using file-based cache (persists between runs)")
    
    weatherAdapter := weather.NewWithCache(log, weatherCache)
    
    app := cli.New(log, weatherAdapter)
    
    if err := app.Run(); err != nil {
        log.Error("Application failed", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}