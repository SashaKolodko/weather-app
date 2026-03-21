package main

import (
    "log"
    "github.com/SashaKolodko/weather-app/internal/pkg/app/cli"
    "github.com/SashaKolodko/weather-app/pkg/logger"
)

func main() {
    log := logger.NewSimpleLogger()
    app := cli.New(log)
    
    if err := app.Run(); err != nil {
        log.Error(fmt.Sprintf("Application failed: %s", err.Error()))
        log.Fatal(err)
    }
}