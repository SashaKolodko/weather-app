#!/bin/bash

# Требуется ImageMagick
# sudo apt install imagemagick

# Создание простой иконки с помощью ImageMagick
convert -size 256x256 xc:none \
    -fill "blue" -draw "rectangle 64,64 192,192" \
    -fill "white" -draw "circle 128,128 128,90" \
    -fill "orange" -draw "circle 128,128 128,120" \
    -fill "white" -pointsize 40 -draw "text 95,140 '°C'" \
    assets/weather-icon.png

# Конвертация в ICO (разные размеры)
convert assets/weather-icon.png \
    -resize 256x256 \
    -define icon:auto-resize=256,128,96,64,48,32,24,16 \
    assets/weather-icon.ico

echo "Icon created at assets/weather-icon.ico"
