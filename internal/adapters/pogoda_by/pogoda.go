package pogodaby

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "weather-app/internal/domain/models"
)

const (
       apiURL = "https://pogoda.by/api/v2/weather-fact?station=26820"
)

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}

type pogodaResponse struct {
    Temp float32 `json:"t"`
}

type PogodaWeather struct {
    logger Logger
    client *http.Client
}

func New(l Logger) *PogodaWeather {
    return &PogodaWeather{
        logger: l,
        client: &http.Client{},
    }
}

func (p *PogodaWeather) GetTemperature(lat, long float64) (models.TempInfo, error) {
    p.logger.Info("Fetching weather data from pogoda.by")
    p.logger.Debug(fmt.Sprintf("Request URL: %s", apiURL))
    
    response, err := p.client.Get(apiURL)
    if err != nil {
        p.logger.Error("can't get data from pogoda.by", err)
        return models.TempInfo{}, err
    }
    defer func() {
        if err := response.Body.Close(); err != nil {
            p.logger.Error("can't close response body", err)
        }
    }()
    
    if response.StatusCode != http.StatusOK {
        errMsg := fmt.Sprintf("pogoda.by API returned non-200 status: %d", response.StatusCode)
        p.logger.Error(errMsg, nil)
        return models.TempInfo{}, fmt.Errorf(errMsg)
    }
    
    var r pogodaResponse
    if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
        p.logger.Error("can't decode JSON", err)
        return models.TempInfo{}, err
    }
    
    p.logger.Info(fmt.Sprintf("Weather data from pogoda.by loaded successfully: %.1f°C", r.Temp))
    
    return models.TempInfo{
        Temp:      r.Temp,
        Humidity:  0,
        WindSpeed: 0,
    }, nil
}
