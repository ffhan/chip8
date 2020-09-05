package wasm

import "syscall/js"

type display struct {
	buffer        []byte
	width, height int
	context       js.Value
}

func NewDisplay(width int, height int) *display {
	canvas := js.Global().Get("document").Call("getElementById", "screen")
	canvas.Set("width", width*10)
	canvas.Set("height", height*10)

	buf := make([]byte, width*height)

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
			idx := d.width*y + x
			d.buffer[idx] = 0
		}
	}
}

func (d *display) Run() error {
	d.Clear()
	return nil
}

func (d *display) Repaint() {
	js.CopyBytesToJS(js.Global().Get("document").Get("vram"), d.buffer)
	js.Global().Get("repaint").Invoke()
}

func (d *display) writeByte(x, y, b byte) bool {
	xb := int(x)
	yb := int(y)
	collision := false
	for i := 0; i < 8; i++ {
		res := (b & 0x80) >> 7
		b <<= 1

		idx := (yb%d.height)*d.width + (xb+i)%d.width

		oldRes := d.buffer[idx]
		d.buffer[idx] ^= res
		newRes := d.buffer[idx]
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
