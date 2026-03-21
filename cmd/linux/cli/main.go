package main

import (
    "os"
    
    "github.com/SashaKolodko/weather-app/internal/pkg/app/cli"
    "github.com/SashaKolodko/weather-app/pkg/logger"
)

func main() {
    // Используем новый логгер из pkg/logger
    log := logger.New()
    
    app := cli.New(log)
    
    err := app.Run()
    if err != nil {
        log.Error("Some error", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}