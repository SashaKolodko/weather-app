package main

import (
    "github.com/SashaKolodko/weather-app/internal/pkg/app/cli"
    "github.com/SashaKolodko/weather-app/pkg/logger"
)

func main() {
    log := logger.NewSimpleLogger()
    app := cli.New(log)
    
    if err := app.Run(); err != nil {
        log.Error("Application failed: " + err.Error())
    }
}