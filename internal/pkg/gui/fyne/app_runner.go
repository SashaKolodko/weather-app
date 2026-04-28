package fyne

import (
    "fyne.io/fyne/v2"
)

type appRunner struct {
    w fyne.Window
}

func NewAR(w fyne.Window) *appRunner {
    return &appRunner{w: w}
}

func (ar *appRunner) Run() {
    ar.w.ShowAndRun()
}
