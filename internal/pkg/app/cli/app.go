package cli

import (
    "fmt"
    
    "weather-app/internal/domain/models"
    "weather-app/pkg/config"
)

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}

type WeatherInfo interface {
    GetTemperature(float64, float64) (models.TempInfo, error)
}

type cliApp struct {
    l    Logger
    wi   WeatherInfo
    conf config.Config
}

func New(l Logger, wi WeatherInfo, c config.Config) *cliApp {
    return &cliApp{
        l:    l,
        wi:   wi,
        conf: c,
    }
}

func (c *cliApp) Run() error {
    c.l.Info("========================================")
    c.l.Info("Starting Weather Application")
    c.l.Info("========================================")
    
    latitude := c.conf.L.Lat
    longitude := c.conf.L.Long
    
    c.l.Info(fmt.Sprintf("Using provider: %s", c.conf.P.Type))
    c.l.Info(fmt.Sprintf("Fetching weather for coordinates: %.4f, %.4f", latitude, longitude))
  
    weather, err := c.wi.GetTemperature(latitude, longitude)
    if err != nil {
        c.l.Error("Failed to get weather data", err)
        return err
    }
    
    c.l.Info("========================================")
    c.l.Info("WEATHER REPORT")
    c.l.Info("========================================")
    
    tempMsg := fmt.Sprintf("🌡️  Температура воздуха: %.1f°C", weather.Temp)
    fmt.Println(tempMsg)
    c.l.Info(tempMsg)
    
    if weather.Humidity != 0 {
        humidityMsg := fmt.Sprintf("💧 Влажность: %.0f%%", weather.Humidity)
        fmt.Println(humidityMsg)
        c.l.Info(humidityMsg)
    }
    
    if weather.WindSpeed != 0 {
        windMsg := fmt.Sprintf("💨 Скорость ветра: %.1f км/ч", weather.WindSpeed)
        fmt.Println(windMsg)
        c.l.Info(windMsg)
    }
    
    c.l.Info("========================================")
    c.l.Info("Weather Application finished successfully")
    c.l.Info("========================================")
    
    return nil
}