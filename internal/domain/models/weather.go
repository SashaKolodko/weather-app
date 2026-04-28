package models

// TempInfo содержит основную информацию о температуре
type TempInfo struct {
    Temp      float32
    Humidity  float32
    WindSpeed float32
}

// WeatherData содержит полную информацию о погоде
type WeatherData struct {
    Temperature float32
    Humidity    float32
    WindSpeed   float32
    Location    string
    Timestamp   string
}