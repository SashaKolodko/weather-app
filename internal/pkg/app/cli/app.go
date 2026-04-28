package cli

import (
    "fmt"
    
    "weather-app/internal/domain/models"
)

// Logger интерфейс логгера
type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}

// WeatherInfo интерфейс для получения информации о погоде
type WeatherInfo interface {
    GetTemperature(float64, float64) models.TempInfo
}

// cliApp структура CLI приложения
type cliApp struct {
    l  Logger
    wi WeatherInfo
}

// New создает новый экземпляр CLI приложения
func New(l Logger, wi WeatherInfo) *cliApp {
    return &cliApp{
        l:  l,
        wi: wi,
    }
}

// Run запускает приложение
func (c *cliApp) Run() error {
    c.l.Info("========================================")
    c.l.Info("Starting Weather Application")
    c.l.Info("========================================")
    
    // Координаты Гродно
    latitude := 53.6688
    longitude := 23.8223
    
    c.l.Info(fmt.Sprintf("Fetching weather for coordinates: %.4f, %.4f", latitude, longitude))
    
    // Получаем данные о погоде
    weather := c.wi.GetTemperature(latitude, longitude)
    
    // Выводим результат
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