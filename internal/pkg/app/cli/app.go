package cli

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "github.com/SashaKolodko/weather-app/pkg/logger"
)

type cliApp struct {
    logger logger.Logger
}

func New(log logger.Logger) *cliApp {
    return &cliApp{
        logger: log,
    }
}

func (c *cliApp) Run() error {
    c.logger.Info("Starting weather application")
    
    type Current struct {
        Temp float32 `json:"temperature_2m"`
    }
    
    type Response struct {
        Curr Current `json:"current"`
    }
    
    var response Response
    
    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m",
        53.6688,
        23.8223,
    )
    
    url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", params)
    
    c.logger.Debug(fmt.Sprintf("Request URL: %s", url))
    
    resp, err := http.Get(url)
    if err != nil {
        c.logger.Error("Failed to get weather data")
        customErr := errors.New("can't get weather data from openmeteo")
        return errors.Join(customErr, err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            c.logger.Error(fmt.Sprintf("Failed to close body: %s", err.Error()))
        }
    }()
    
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        c.logger.Error("Failed to read response data")
        customErr := errors.New("can't read data from response")
        return errors.Join(customErr, err)
    }
    
    if err := json.Unmarshal(data, &response); err != nil {
        c.logger.Error("Failed to unmarshal JSON")
        customErr := errors.New("can't unmarshal data from response")
        return errors.Join(customErr, err)
    }
    
    result := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", response.Curr.Temp)
    c.logger.Info(result)
    fmt.Println(result)
    
    return nil
}