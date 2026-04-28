package fyne

import (
    "fmt"
    
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/widget"
    guisettings "weather-app/internal/domain/gui_settings"
)

type window struct {
    w      fyne.Window
    label  *widget.Label
}

func NewW(w fyne.Window) *window {
    return &window{
        w: w,
    }
}

func (win *window) Resize(ws guisettings.WindowSize) error {
    if !ws.IsFull() {
        win.w.Resize(fyne.NewSize(float32(ws.Width()), float32(ws.Height())))
    } else {
        win.w.SetFullScreen(true)
    }
    return nil
}

func (win *window) UpdateTemperature(t float32) error {
    if win.label != nil {
        win.label.SetText(fmt.Sprintf("Температура: %.1f°C", t))
    }
    return nil
}

func (win *window) SetTemperatureWidget(tw guisettings.TextWidget) error {
    if widget, ok := tw.Render().(*widget.Label); ok {
        win.label = widget
        win.w.SetContent(widget)
    }
    return nil
}

func (win *window) Render() error {
    win.w.Show()
    return nil
}
