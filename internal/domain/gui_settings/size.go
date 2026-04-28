package guisettings

type WindowSize struct {
    width  int
    height int
}

func NewWS(w, h int) WindowSize {
    return WindowSize{width: w, height: h}
}

func (ws WindowSize) IsFull() bool {
    return ws.width == 0 && ws.height == 0
}

func (ws WindowSize) Width() int {
    return ws.width
}

func (ws WindowSize) Height() int {
    return ws.height
}
