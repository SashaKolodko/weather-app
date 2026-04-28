package main

import (
    "os"
    
    "weather-app/internal/adapters/weather"
    "weather-app/internal/pkg/app/cli"
    "weather-app/internal/pkg/flags"
    "weather-app/pkg/cache"
    "weather-app/pkg/config"
    "weather-app/pkg/logger"
)

func main() {
    args := flags.Parse()
    
    r, err := os.Open(args.Path)
    if err != nil {
        panic(err)
    }
    defer r.Close()
    
    cfg, err := config.Parse(r)
    if err != nil {
        panic(err)
    }
    
    log := logger.New()
    
    weatherCache, err := cache.NewFileCache("")
    if err != nil {
        log.Error("Failed to create cache", err)
        os.Exit(1)
    }
    
    log.Info("Using file-based cache (persists between runs)")
    
    wi := getProvider(cfg, log, weatherCache)
    
    app := cli.New(log, wi, cfg)
    
    if err := app.Run(); err != nil {
        log.Error("Application failed", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}

func getProvider(cfg config.Config, log *logger.Logger, cache *cache.FileCache) cli.WeatherInfo {
    var wi cli.WeatherInfo
    
    switch cfg.P.Type {
    case "open-meteo":
        wi = weather.NewWithCache(log, cache)
    default:
        wi = weather.NewWithCache(log, cache)
    }
    
    return wi
}