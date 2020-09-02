package main

import (
	"bytes"
	"chip8"
	"syscall/js"
)

type speaker struct {
}

func (s *speaker) Play() {
	js.Global().Get("playMusic").Invoke()
}

func (s *speaker) Stop() {
	js.Global().Get("stopMusic").Invoke()
}

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

func (s *screen) Run() error {
	s.Clear()
	return nil
}

func (s *screen) writeByte(x, y, b byte) bool {
	xb := int(x)
	yb := int(y)
	collision := false
	for i := 0; i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1
		oldRes := s.buffer[yb%s.height][(xb+i)%s.width]
		s.buffer[yb%s.height][(xb+i)%s.width] ^= res
		newRes := s.buffer[yb%s.height][(xb+i)%s.width]
		if oldRes == 1 && newRes == 0 {
			collision = true
		}
		if oldRes != newRes {
			s.write((xb+i)%s.width, yb%s.height, newRes)
		}
	}
	return collision
}

func (s *screen) write(x int, y int, value byte) {
	if value == 1 {
		s.context.Set("fillStyle", "white")
	} else {
		s.context.Set("fillStyle", "black")
	}
	s.context.Call("fillRect", x*10, y*10, 10, 10)
}

func (s *screen) Write(x, y byte, bytes []byte) bool {
	collision := false
	for i := range bytes {
		if s.writeByte(x, y+byte(i), bytes[i]) {
			collision = true
		}
	}
	return collision
}

func run(this js.Value, i []js.Value) interface{} {
	rom := make([]byte, 0xFFF)
	n := js.CopyBytesToGo(rom, js.Global().Get("document").Get("buffer"))

	rom = rom[:n]
	romBuf := bytes.NewBuffer(rom)

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

	keyboard := chip8.NewDefaultKeyboard()

	js.Global().Set("down", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyboard.KeyDown(rune(args[0].String()[0]))
		return nil
	}))
	js.Global().Set("up", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyboard.KeyUp(rune(args[0].String()[0]))
		return nil
	}))

	chip8.Run(romBuf, chip8.Chip8, s, keyboard, &speaker{})
	return nil
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")

	js.Global().Set("run", js.FuncOf(run))

	<-c
}
