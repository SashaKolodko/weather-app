; Inno Setup Script for Weather App
; Compile with: ISCC weather-app.iss

[Setup]
; Основные настройки
AppId={{WEATHER-APP-GUID-1234-5678-90AB-CDEF12345678}
AppName=Weather App
AppVersion=1.0.0
AppPublisher=Sasha Kolodko
AppPublisherURL=https://github.com/SashaKolodko/weather-app
AppSupportURL=https://github.com/SashaKolodko/weather-app
AppUpdatesURL=https://github.com/SashaKolodko/weather-app
DefaultDirName={autopf}\WeatherApp
DefaultGroupName=Weather App
AllowNoIcons=yes
LicenseFile=..\LICENSE
OutputDir=..\build\installer
OutputBaseFilename=WeatherApp-Setup
SetupIconFile=..\assets\weather-icon.ico
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=lowest
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64

; Требования
MinVersion=10.0.14393

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "russian"; MessagesFile: "compiler:Languages\Russian.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked
Name: "startmenuicon"; Description: "{cm:CreateStartMenuIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Основные исполняемые файлы
Source: "..\build\windows\weather-cli.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\build\windows\weather-gui.exe"; DestDir: "{app}"; Flags: ignoreversion

; Конфигурация
Source: "..\build\windows\config\*.yaml"; DestDir: "{app}\config"; Flags: ignoreversion recursesubdirs createallsubdirs

; Документация
Source: "..\README.md"; DestDir: "{app}"; Flags: isreadme

; Создание пустых директорий для данных
; (будет создано при первом запуске)

[Icons]
Name: "{group}\Weather App (GUI)"; Filename: "{app}\weather-gui.exe"; Tasks: startmenuicon
Name: "{group}\Weather App (CLI)"; Filename: "{app}\weather-cli.exe"; Tasks: startmenuicon
Name: "{group}\Uninstall Weather App"; Filename: "{uninstallexe}"
Name: "{autodesktop}\Weather App"; Filename: "{app}\weather-gui.exe"; Tasks: desktopicon

[Run]
; Запуск GUI после установки (опционально)
Filename: "{app}\weather-gui.exe"; Description: "{cm:LaunchProgram,Weather App}"; Flags: postinstall nowait skipifsilent unchecked

; Создание ярлыка в PATH (для CLI)
[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Flags: preservestringtype uninsdeletevalue

[UninstallDelete]
; Удаление кэша и настроек
Type: filesandordirs; Name: "{userappdata}\WeatherApp"

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
var
  ResultCode: Integer;
begin
  if CurStep = ssPostInstall then
  begin
    // Создание директорий для кэша
    if not DirExists(ExpandConstant('{userappdata}\WeatherApp')) then
      CreateDir(ExpandConstant('{userappdata}\WeatherApp'));
    
    // Копирование примера конфигурации
    if FileExists(ExpandConstant('{app}\config\config.yaml.example')) then
      FileCopy(ExpandConstant('{app}\config\config.yaml.example'), 
               ExpandConstant('{userappdata}\WeatherApp\config.yaml'), False);
  end;
end;

// Добавление в PATH
function GetEnv(const Name: string): string;
begin
  Result := GetEnvVar(Name);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
  begin
    // Предложение удалить данные пользователя
    if MsgBox('Do you want to remove user data (cache and settings)?', 
              mbConfirmation, MB_YESNO) = IDYES then
    begin
      DelTree(ExpandConstant('{userappdata}\WeatherApp'), True, True, True);
    end;
  end;
end;
