#!/bin/bash

set -e

echo "=========================================="
echo "Building Weather App for Windows"
echo "=========================================="

# Создание директории
mkdir -p build/windows-release

# Сборка Windows бинарников
echo "Building Windows binaries..."

# Сборка CLI
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" \
    -o build/windows-release/weather-cli.exe \
    ./cmd/linux/cli/main.go

# Сборка GUI (требует CGO)
echo "Building GUI (may take a moment)..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
    CC=x86_64-w64-mingw32-gcc \
    go build -ldflags="-s -w" \
    -o build/windows-release/weather-gui.exe \
    ./cmd/linux/gui/main.go 2>/dev/null || \
    echo "GUI build skipped (requires mingw), building only CLI..."

# Копирование конфигов
mkdir -p build/windows-release/config
cp config/*.yaml build/windows-release/config/ 2>/dev/null || true

# Копирование README
cp README.md build/windows-release/

# Создание запускных скриптов
cat > build/windows-release/run-gui.bat << 'END'
@echo off
title Weather App
echo Starting Weather App GUI...
start weather-gui.exe
exit
END

cat > build/windows-release/run-cli.bat << 'END'
@echo off
title Weather App CLI
echo Weather App Command Line Interface
echo ================================
weather-cli.exe
echo.
pause
END

# Создание скрипта установки
cat > build/windows-release/install.bat << 'END'
@echo off
echo Installing Weather App...
echo.

:: Создание директории Program Files
if not exist "%ProgramFiles%\WeatherApp" (
    mkdir "%ProgramFiles%\WeatherApp"
)

:: Копирование файлов
xcopy /E /I /Y "%~dp0*" "%ProgramFiles%\WeatherApp\"

:: Создание ярлыка на рабочем столе
echo Creating desktop shortcut...
powershell -Command "$WS = New-Object -ComObject WScript.Shell; $SC = $WS.CreateShortcut('%USERPROFILE%\Desktop\Weather App.lnk'); $SC.TargetPath = '%ProgramFiles%\WeatherApp\weather-gui.exe'; $SC.Save()"

:: Добавление в PATH
echo Adding to PATH...
setx Path "%Path%;%ProgramFiles%\WeatherApp"

echo.
echo Installation complete!
echo You can now run Weather App from Start Menu or Desktop
pause
END

# Создание скрипта удаления
cat > build/windows-release/uninstall.bat << 'END'
@echo off
echo Uninstalling Weather App...
echo.

:: Удаление директории
rmdir /S /Q "%ProgramFiles%\WeatherApp"

:: Удаление ярлыка
del "%USERPROFILE%\Desktop\Weather App.lnk"

echo.
echo Uninstallation complete!
pause
END

# Создание ZIP архива
echo "Creating ZIP archive..."
cd build/windows-release
zip -r ../WeatherApp-Windows.zip . 2>/dev/null || {
    # Если zip не установлен, используем tar
    tar -czf ../WeatherApp-Windows.tar.gz . 2>/dev/null
    echo "Created WeatherApp-Windows.tar.gz"
}
cd ../..

echo "=========================================="
echo "Build complete!"
echo "=========================================="
echo "Output files:"
ls -lh build/WeatherApp-Windows.* 2>/dev/null || echo "  build/windows-release/"
echo "=========================================="
