package wasm

import "syscall/js"

type display struct {
	buffer        [][]byte
	width, height int
	context       js.Value
}

func NewDisplay(width int, height int) *display {
	canvas := js.Global().Get("document").Call("getElementById", "screen")
	canvas.Set("width", width*10)
	canvas.Set("height", height*10)

	buf := make([][]byte, height)
	for i := range buf {
		buf[i] = make([]byte, width)
	}

	d := &display{
		buffer:  buf,
		width:   width,
		height:  height,
		context: canvas.Call("getContext", "2d"),
	}
	return d
}

func (d *display) Clear() {
	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			d.buffer[y][x] = 0
		}
	}
}

func (d *display) Run() error {
	d.Clear()
	return nil
}

func (d *display) Repaint() {
	for y := range d.buffer {
		for x := range d.buffer[y] {
			d.write(x, y, d.buffer[y][x])
		}
	}
}

func (d *display) writeByte(x, y, b byte) bool {
	xb := int(x)
	yb := int(y)
	collision := false
	for i := 0; i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1
		oldRes := d.buffer[yb%d.height][(xb+i)%d.width]
		d.buffer[yb%d.height][(xb+i)%d.width] ^= res
		newRes := d.buffer[yb%d.height][(xb+i)%d.width]
		if oldRes == 1 && newRes == 0 {
			collision = true
		}
	}
	return collision
}

func (d *display) write(x int, y int, value byte) {
	if value == 1 {
		d.context.Set("fillStyle", "white")
	} else {
		d.context.Set("fillStyle", "black")
	}
	d.context.Call("fillRect", x*10, y*10, 10, 10)
}

func (d *display) Write(x, y byte, bytes []byte) bool {
	collision := false
	for i := range bytes {
		if d.writeByte(x, y+byte(i), bytes[i]) {
			collision = true
		}
	}
	return collision
}
