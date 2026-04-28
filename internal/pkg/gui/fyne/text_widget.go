package fyne

import (
    "fyne.io/fyne/v2/widget"
)

type TextWidget struct {
    l *widget.Label
}

func NewTW(label *widget.Label) *TextWidget {
    return &TextWidget{l: label}
}

func (tw *TextWidget) Render() any {
    tw.l.Show()
    return tw.l
}

func (tw *TextWidget) SetText(text string) {
    tw.l.SetText(text)
}
