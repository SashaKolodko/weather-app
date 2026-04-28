package weather

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "weather-app/internal/domain/models"
    "weather-app/pkg/cache"
)

const (
    apiURL         = "https://api.open-meteo.com/v1/forecast"
    defaultTimeout = 30 * time.Second
    cacheTTL       = 10 * time.Minute
)


type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}


type current struct {
    Temp      float32 `json:"temperature_2m"`
    Humidity  float32 `json:"relative_humidity_2m"`
    WindSpeed float32 `json:"wind_speed_10m"`
}


type apiResponse struct {
    Current current `json:"current"`
}


type OpenMeteoWeather struct {
    logger   Logger
    client   *http.Client
    cache    cache.Cache
    isLoaded bool
    lastData *models.WeatherData
}


func New(l Logger) *OpenMeteoWeather {
    return &OpenMeteoWeather{
        logger:   l,
        client:   &http.Client{Timeout: defaultTimeout},
        isLoaded: false,
    }
}


func NewWithCache(l Logger, c cache.Cache) *OpenMeteoWeather {
    return &OpenMeteoWeather{
        logger:   l,
        client:   &http.Client{Timeout: defaultTimeout},
        cache:    c,
        isLoaded: false,
    }
}


func (w *OpenMeteoWeather) generateCacheKey(lat, long float64) string {
    data := fmt.Sprintf("%f:%f", lat, long)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}


func (w *OpenMeteoWeather) fetchWeatherData(lat, long float64) error {
    w.logger.Debug(fmt.Sprintf("Fetching weather data for coordinates: %.4f, %.4f", lat, long))
    
    
    if w.cache != nil {
        cacheKey := w.generateCacheKey(lat, long)
        if cachedData, found := w.cache.Get(cacheKey); found {
            w.logger.Info("✓ Weather data found in cache")
            
            var response apiResponse
            if err := json.Unmarshal(cachedData, &response); err == nil {
                w.lastData = &models.WeatherData{
                    Temperature: response.Current.Temp,
                    Humidity:    response.Current.Humidity,
                    WindSpeed:   response.Current.WindSpeed,
                    Timestamp:   time.Now().Format(time.RFC3339),
                }
                w.isLoaded = true
                return nil
            }
        }
        w.logger.Info("✗ Cache miss, fetching from API")
    }
    
    
    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m&timezone=auto",
        lat, long,
    )
    url := fmt.Sprintf("%s?%s", apiURL, params)
    
    w.logger.Debug(fmt.Sprintf("Request URL: %s", url))
    
    
    resp, err := w.client.Get(url)
    if err != nil {
        w.logger.Error("Failed to get weather data", err)
        return errors.New("can't get weather data from openmeteo")
    }
    defer resp.Body.Close()
    
    
    if resp.StatusCode != http.StatusOK {
        errMsg := fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode)
        w.logger.Error(errMsg, nil)
        return errors.New(errMsg)
    }
    
    
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        w.logger.Error("Failed to read response data", err)
        return errors.New("can't read data from response")
    }
    
    w.logger.Debug(fmt.Sprintf("Data fetched successfully, size: %d bytes", len(data)))
    
    
    if w.cache != nil {
        cacheKey := w.generateCacheKey(lat, long)
        w.cache.Set(cacheKey, data, cacheTTL)
        w.logger.Info(fmt.Sprintf("✓ Weather data saved to cache (TTL: %v)", cacheTTL))
    }
    
    
    var response apiResponse
    if err := json.Unmarshal(data, &response); err != nil {
        w.logger.Error("Failed to unmarshal JSON", err)
        return errors.New("can't unmarshal data from response")
    }
    

    w.lastData = &models.WeatherData{
        Temperature: response.Current.Temp,
        Humidity:    response.Current.Humidity,
        WindSpeed:   response.Current.WindSpeed,
        Timestamp:   time.Now().Format(time.RFC3339),
    }
    w.isLoaded = true
    
    w.logger.Info("Weather data loaded successfully")
    w.logger.Debug(fmt.Sprintf("Temperature: %.1f°C, Humidity: %.0f%%, Wind: %.1f km/h",
        w.lastData.Temperature, w.lastData.Humidity, w.lastData.WindSpeed))
    
    return nil
}

func (w *OpenMeteoWeather) GetTemperature(lat, long float64) (models.TempInfo, error) {
    if err := w.fetchWeatherData(lat, long); err != nil {
        w.logger.Error("Failed to get temperature", err)
        return models.TempInfo{}, err
    }
    
    return models.TempInfo{
        Temp:      w.lastData.Temperature,
        Humidity:  w.lastData.Humidity,
        WindSpeed: w.lastData.WindSpeed,
    }, nil
}


func (w *OpenMeteoWeather) GetWeatherData(lat, long float64) (*models.WeatherData, error) {
    if err := w.fetchWeatherData(lat, long); err != nil {
        return nil, err
    }
    return w.lastData, nil
}