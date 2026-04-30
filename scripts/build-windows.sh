#!/bin/bash

# Создание директории для сборки
mkdir -p build/windows

# Переменные
APP_NAME="WeatherApp"
CLI_BINARY="weather-cli.exe"
GUI_BINARY="weather-gui.exe"
VERSION="1.0.0"

echo "Building Windows binaries..."

# Сборка CLI версии для Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" \
    -o build/windows/$CLI_BINARY \
    ./cmd/linux/cli/main.go

# Сборка GUI версии для Windows (Fyne требует CGO)
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
    go build -ldflags="-s -w" \
    -o build/windows/$GUI_BINARY \
    ./cmd/linux/gui/main.go

# Копирование конфигурационных файлов
cp -r config build/windows/
cp README.md build/windows/

# Создание директории для данных
mkdir -p build/windows/data

# Создание пустого файла кэша
mkdir -p build/windows/cache

echo "Windows binaries built successfully in build/windows/"
ls -la build/windows/
