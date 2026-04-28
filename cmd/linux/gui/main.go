package main

import (
    "os"
    
    "weather-app/internal/pkg/app/gui"
    "weather-app/internal/pkg/flags"
    "weather-app/internal/pkg/gui/fyne"
    "weather-app/internal/adapters/weather"
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
    
    // Парсим конфигурацию
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
    
    var weatherProvider gui.WeatherInfo
    switch cfg.P.Type {
    case "open-meteo":
        log.Info("Using Open-Meteo weather provider")
        weatherProvider = weather.NewWithCache(log, weatherCache)
    default:
        log.Info("Using default Open-Meteo weather provider")
        weatherProvider = weather.NewWithCache(log, weatherCache)
    }
    
    guiProvider := fyne.NewP()
    
    app := gui.New(log, guiProvider, weatherProvider, cfg)
    
    if err := app.Run(); err != nil {
        log.Error("GUI Application failed", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}
