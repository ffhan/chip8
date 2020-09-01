package chip8

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

const (
	pixelFactor = 10
)

type Display interface {
	Clear()
	Write(x, y byte, bytes []byte) bool
}

type myRaster struct {
	widget.BaseWidget
	width, height int
	buffer        [][]byte
	raster        *canvas.Raster
	keyboard      Keyboard
}

type myRenderer struct {
	render *canvas.Raster
}

func (m *myRenderer) Layout(size fyne.Size) {
	m.render.Resize(size)
}

func (m *myRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (m *myRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{m.render}
}

func (m *myRenderer) MinSize() fyne.Size {
	return m.render.MinSize()
}

func (m *myRenderer) Refresh() {
	m.render.Refresh()
}

func (m *myRenderer) Destroy() {
}

func (m *myRaster) CreateRenderer() fyne.WidgetRenderer {
	return &myRenderer{render: m.raster}
}

func (m *myRaster) FocusGained() {
}

func (m *myRaster) FocusLost() {
}

func (m *myRaster) Focused() bool {
	return true
}

func (m *myRaster) TypedRune(r rune) {
}

func (m *myRaster) TypedKey(event *fyne.KeyEvent) {
}

func (m *myRaster) KeyDown(event *fyne.KeyEvent) {
	m.keyboard.KeyDown(rune(event.Name[0]))
}

func (m *myRaster) KeyUp(event *fyne.KeyEvent) {
	m.keyboard.KeyUp(rune(event.Name[0]))
}

type guiDisplay struct {
	buffer        [][]byte // monochrome!
	width, height int
	raster        *myRaster
	window        fyne.Window
	keyboard      Keyboard
}

func NewGuiDisplay(width, height int, keyboard Keyboard) *guiDisplay {
	buffer := make([][]byte, height)
	for i := range buffer {
		buffer[i] = make([]byte, width)
	}
	a := app.New()
	w := a.NewWindow("Chip-8 emulator")
	size := fyne.NewSize(width*pixelFactor, height*pixelFactor)
	w.Resize(size)
	w.CenterOnScreen()
	w.SetFixedSize(true)

	c := &guiDisplay{
		buffer:   buffer,
		width:    width,
		height:   height,
		window:   w,
		keyboard: keyboard,
	}

	raster := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if x >= width*pixelFactor || y >= height*pixelFactor {
			return color.White
		}
		x /= pixelFactor
		y /= pixelFactor
		if c.buffer[y][x] > 0 {
			return color.White
		}
		return color.Black
	})

	cRaster := &myRaster{
		BaseWidget: widget.BaseWidget{},
		width:      width,
		height:     height,
		buffer:     buffer,
		raster:     raster,
		keyboard:   keyboard,
	}
	cRaster.ExtendBaseWidget(cRaster)
	w.SetContent(cRaster)
	c.raster = cRaster
	return c
}

func (g *guiDisplay) GetCanvas() fyne.Canvas {
	return g.window.Canvas()
}

func (g *guiDisplay) Run() {
	g.window.ShowAndRun()
}

func NewDefaultGuiDisplay(keyboard Keyboard) *guiDisplay {
	return NewGuiDisplay(64, 32, keyboard)
}

func (g *guiDisplay) writeByte(x, y, b byte) bool {
	xb := int(x)
	yb := int(y)
	collision := false
	for i := 0; i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1
		oldRes := g.buffer[yb%g.height][(xb+i)%g.width]
		g.buffer[yb%g.height][(xb+i)%g.width] ^= res
		if oldRes == 1 && g.buffer[yb%g.height][(xb+i)%g.width] == 0 {
			collision = true
		}
	}
	return collision
}

func (g *guiDisplay) Write(x, y byte, bytes []byte) bool {
	collision := false
	for i := range bytes {
		if g.writeByte(x, y+byte(i), bytes[i]) {
			collision = true
		}
	}
	g.updateScreen()
	return collision
}

func (g *guiDisplay) Clear() {
	for i := range g.buffer {
		for j := range g.buffer[i] {
			g.buffer[i][j] = 0
		}
	}
	g.updateScreen()
}

func (g *guiDisplay) updateScreen() {
	g.raster.Refresh()
}
