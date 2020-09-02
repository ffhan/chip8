package main

import (
	"syscall/js"
)

type screen struct {
	buffer        [][]byte
	width, height int
	context       js.Value
}

func (s *screen) Clear() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.write(x, y, 0)
		}
	}
}

func (g *screen) writeByte(x, y, b byte) bool {
	xb := int(x)
	yb := int(y)
	collision := false
	for i := 0; i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1
		oldRes := g.buffer[yb%g.height][(xb+i)%g.width]
		g.buffer[yb%g.height][(xb+i)%g.width] ^= res
		newRes := g.buffer[yb%g.height][(xb+i)%g.width]
		if oldRes == 1 && newRes == 0 {
			collision = true
		}
		if oldRes != newRes {
			g.write((xb+i)%g.width, yb%g.height, newRes)
		}
	}
	return collision
}

func (g *screen) write(x int, y int, value byte) {
	if value == 1 {
		g.context.Set("fillStyle", "white")
	} else {
		g.context.Set("fillStyle", "black")
	}
	g.context.Call("fillRect", x*10, y*10, 10, 10)
}

func (g *screen) Write(x, y byte, bytes []byte) bool {
	collision := false
	for i := range bytes {
		if g.writeByte(x, y+byte(i), bytes[i]) {
			collision = true
		}
	}
	return collision
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")

	width := 64
	height := 32

	buf := make([][]byte, height)
	for i := range buf {
		buf[i] = make([]byte, width)
	}

	canvas := js.Global().Get("document").Call("getElementById", "screen")
	canvas.Set("width", width*10)
	canvas.Set("height", height*10)

	s := &screen{
		buffer:  buf,
		width:   width,
		height:  height,
		context: canvas.Call("getContext", "2d"),
	}

	SpriteA := []byte{0xF0, 0x90, 0xF0, 0x90, 0x90}
	SpriteC := []byte{0xF0, 0x80, 0x80, 0x80, 0xF0}

	s.Clear()
	s.Write(10, 31, SpriteA)
	s.Write(40, 10, SpriteC)

	<-c
}
