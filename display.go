package chip8

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"image/color"
)

const (
	pixelFactor = 10
)

type Display interface {
	Clear()
	Write(x, y byte, bytes []byte)
	Get(pointer uint16) byte
}

type guiDisplay struct {
	buffer        [][]byte // monochrome!
	width, height int
	raster        *canvas.Raster
	window        fyne.Window
}

func NewGuiDisplay(width, height int) *guiDisplay {
	buffer := make([][]byte, height)
	for i := range buffer {
		buffer[i] = make([]byte, width)
	}
	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(width*pixelFactor, height*pixelFactor))
	w.CenterOnScreen()
	w.SetFixedSize(true)

	c := &guiDisplay{
		buffer: buffer,
		width:  width,
		height: height,
		window: w,
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
	w.SetContent(raster)
	c.raster = raster
	return c
}

func (g *guiDisplay) Run() {
	g.window.ShowAndRun()
}

func NewDefaultGuiDisplay() *guiDisplay {
	return NewGuiDisplay(64, 32)
}

func (g *guiDisplay) writeByte(x, y, b byte) uint16 {
	for i := byte(0); i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1
		g.buffer[y][x+i] = res
	}
	return 8
}

func (g *guiDisplay) Write(x, y byte, bytes []byte) {
	for i := range bytes {
		g.writeByte(x, y+byte(i), bytes[i])
	}
	g.updateScreen()
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
