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
    apiURL        = "https://api.open-meteo.com/v1/forecast"
    defaultTimeout = 30 * time.Second
    cacheTTL      = 10 * time.Minute
)

// Logger интерфейс для логирования
type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}

// current представляет текущую погоду из API
type current struct {
    Temp      float32 `json:"temperature_2m"`
    Humidity  float32 `json:"relative_humidity_2m"`
    WindSpeed float32 `json:"wind_speed_10m"`
}

// apiResponse представляет ответ от API
type apiResponse struct {
    Current current `json:"current"`
}

// OpenMeteoWeather адаптер для API Open-Meteo
type OpenMeteoWeather struct {
    logger   Logger
    client   *http.Client
    cache    cache.Cache
    isLoaded bool
    lastData *models.WeatherData
}

// New создает новый экземпляр адаптера погоды без кэша
func New(l Logger) *OpenMeteoWeather {
    return &OpenMeteoWeather{
        logger:   l,
        client:   &http.Client{Timeout: defaultTimeout},
        isLoaded: false,
    }
}

// NewWithCache создает новый экземпляр адаптера погоды с кэшем
func NewWithCache(l Logger, c cache.Cache) *OpenMeteoWeather {
    return &OpenMeteoWeather{
        logger:   l,
        client:   &http.Client{Timeout: defaultTimeout},
        cache:    c,
        isLoaded: false,
    }
}

// generateCacheKey генерирует ключ для кэша
func (w *OpenMeteoWeather) generateCacheKey(lat, long float64) string {
    data := fmt.Sprintf("%f:%f", lat, long)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

// fetchWeatherData получает данные о погоде из API
func (w *OpenMeteoWeather) fetchWeatherData(lat, long float64) error {
    w.logger.Debug(fmt.Sprintf("Fetching weather data for coordinates: %.4f, %.4f", lat, long))
    
    // Проверяем кэш
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
    
    // Формируем параметры запроса
    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m&timezone=auto",
        lat, long,
    )
    url := fmt.Sprintf("%s?%s", apiURL, params)
    
    w.logger.Debug(fmt.Sprintf("Request URL: %s", url))
    
    // Выполняем запрос
    resp, err := w.client.Get(url)
    if err != nil {
        w.logger.Error("Failed to get weather data", err)
        return errors.New("can't get weather data from openmeteo")
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            w.logger.Error("Failed to close response body", err)
        }
    }()
    
    // Проверяем статус ответа
    if resp.StatusCode != http.StatusOK {
        errMsg := fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode)
        w.logger.Error(errMsg, nil)
        return errors.New(errMsg)
    }
    
    // Читаем данные
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        w.logger.Error("Failed to read response data", err)
        return errors.New("can't read data from response")
    }
    
    w.logger.Debug(fmt.Sprintf("Data fetched successfully, size: %d bytes", len(data)))
    
    // Сохраняем в кэш
    if w.cache != nil {
        cacheKey := w.generateCacheKey(lat, long)
        w.cache.Set(cacheKey, data, cacheTTL)
        w.logger.Info(fmt.Sprintf("✓ Weather data saved to cache (TTL: %v)", cacheTTL))
    }
    
    // Парсим JSON
    var response apiResponse
    if err := json.Unmarshal(data, &response); err != nil {
        w.logger.Error("Failed to unmarshal JSON", err)
        return errors.New("can't unmarshal data from response")
    }
    
    // Сохраняем данные
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

// GetTemperature возвращает температуру для указанных координат
func (w *OpenMeteoWeather) GetTemperature(lat, long float64) models.TempInfo {
    // Если данные не загружены, загружаем
    if !w.isLoaded {
        if err := w.fetchWeatherData(lat, long); err != nil {
            w.logger.Error("Failed to get temperature", err)
            return models.TempInfo{Temp: 0, Humidity: 0, WindSpeed: 0}
        }
    }
    
    return models.TempInfo{
        Temp:      w.lastData.Temperature,
        Humidity:  w.lastData.Humidity,
        WindSpeed: w.lastData.WindSpeed,
    }
}

// GetWeatherData возвращает полные данные о погоде
func (w *OpenMeteoWeather) GetWeatherData(lat, long float64) (*models.WeatherData, error) {
    if err := w.fetchWeatherData(lat, long); err != nil {
        return nil, err
    }
    return w.lastData, nil
}