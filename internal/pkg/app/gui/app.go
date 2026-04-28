package gui

import (
    "fmt"
    
    guisettings "weather-app/internal/domain/gui_settings"
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

type GUIParams interface {
    CreateWindow(name string, size guisettings.WindowSize) (guisettings.Window, error)
    GetAppRunner() guisettings.AppRunner
    GetTextWidget(text string) guisettings.TextWidget
}

type GUIApp struct {
    logger   Logger
    params   GUIParams
    weather  WeatherInfo
    config   config.Config
    window   guisettings.Window
    textWidget guisettings.TextWidget
}

func New(l Logger, p GUIParams, w WeatherInfo, c config.Config) *GUIApp {
    return &GUIApp{
        logger:   l,
        params:   p,
        weather:  w,
        config:   c,
    }
}

func (g *GUIApp) Run() error {
    g.logger.Info("Starting GUI Weather Application")
    
    
    windowSize := guisettings.NewWS(400, 300)
    window, err := g.params.CreateWindow("Weather App", windowSize)
    if err != nil {
        g.logger.Error("Failed to create window", err)
        return err
    }
    g.window = window
    
    g.textWidget = g.params.GetTextWidget("Loading weather data...")
    if err := g.window.SetTemperatureWidget(g.textWidget); err != nil {
        g.logger.Error("Failed to set temperature widget", err)
        return err
    }
    
    weather, err := g.weather.GetTemperature(g.config.L.Lat, g.config.L.Long)
    if err != nil {
        g.logger.Error("Failed to get weather data", err)
        g.textWidget.SetText(fmt.Sprintf("Error: %v", err))
    } else {
        tempMsg := fmt.Sprintf("Температура: %.1f°C\nВлажность: %.0f%%\nВетер: %.1f км/ч",
            weather.Temp, weather.Humidity, weather.WindSpeed)
        g.textWidget.SetText(tempMsg)
        g.logger.Info(fmt.Sprintf("Weather data loaded: %.1f°C", weather.Temp))
    }
    
    if err := g.window.Render(); err != nil {
        g.logger.Error("Failed to render window", err)
        return err
    }
    
    runner := g.params.GetAppRunner()
    runner.Run()
    
    return nil
}
