package main

import (
	"bytes"
	"chip8"
	"chip8/wasm"
	"syscall/js"
)

func run(this js.Value, i []js.Value) interface{} {
	rom := make([]byte, 0xFFF)
	n := js.CopyBytesToGo(rom, js.Global().Get("document").Get("buffer"))

	rom = rom[:n]
	romBuf := bytes.NewBuffer(rom)

	width := 64
	height := 32

	keyboard := chip8.NewDefaultKeyboard()
	registerKeyboard(keyboard)

	stopFunc := chip8.Run(romBuf, chip8.Chip8, wasm.NewDisplay(width, height), keyboard, wasm.NewSpeaker())
	js.Global().Set("stopCpu", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		stopFunc()
		return nil
	}))
	return nil
}

func registerKeyboard(keyboard chip8.Keyboard) {
	js.Global().Set("down", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyboard.KeyDown(rune(args[0].String()[0]))
		return nil
	}))
	js.Global().Set("up", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyboard.KeyUp(rune(args[0].String()[0]))
		return nil
	}))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")

	js.Global().Set("run", js.FuncOf(run))

	<-c
}
