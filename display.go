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
	Write(x, y byte, bytes []byte) bool
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

func (g *guiDisplay) writeByte(x, y, b byte) bool {
	xb := int(x)
	yb := int(y)
	collision := false
	for i := 0; i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1
		oldRes := g.buffer[yb%g.height][(xb+i)%g.width]
		g.buffer[yb%g.height][(xb+i)%g.width] ^= res
		if oldRes == 1 && g.buffer[yb%g.height][(xb+i)%g.width] == 1 {
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
