package providers

import (
    "weather-app/internal/adapters/weather"
    "weather-app/internal/pkg/app/cli"
    "weather-app/pkg/config"
)

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}

func GetProvider(cfg config.Config, l Logger) cli.WeatherInfo {
    var wi cli.WeatherInfo
    
    switch cfg.P.Type {
    case "open-meteo":
        l.Info("Using Open-Meteo weather provider")
        wi = weather.New(l)
    case "pogoda":
        l.Info("Using Pogoda.by weather provider")
        wi = weather.New(l) 
    default:
        l.Info("Unknown provider, defaulting to Open-Meteo")
        wi = weather.New(l)
    }
    
    return wi
}
