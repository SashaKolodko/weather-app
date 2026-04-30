# Windows installation script wrapper
param(
    [switch]$NoGUI
)

Write-Host "Weather App Installation" -ForegroundColor Cyan
Write-Host "========================" -ForegroundColor Cyan

$installPath = "$env:ProgramFiles\WeatherApp"

# Создание директории
if (!(Test-Path $installPath)) {
    New-Item -ItemType Directory -Path $installPath -Force
}

# Копирование файлов
Copy-Item ".\build\windows\*" -Destination $installPath -Recurse -Force

# Добавление в PATH
$envPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($envPath -notlike "*$installPath*") {
    [Environment]::SetEnvironmentVariable("Path", "$envPath;$installPath", "User")
}

# Создание ярлыка
$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\Weather App.lnk")
$Shortcut.TargetPath = "$installPath\weather-gui.exe"
$Shortcut.Save()

Write-Host "Installation complete!" -ForegroundColor Green
Write-Host "You can now run Weather App from Start Menu or Desktop" -ForegroundColor Yellow

# Запуск
if (!$NoGUI) {
    Start-Process "$installPath\weather-gui.exe"
}
