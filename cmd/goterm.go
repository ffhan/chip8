package main

import (
	"chip8"
)

func main() {
	display := chip8.NewDefaultGuiDisplay()
	//display.Write(10,32)
	//display.Write(128,32)
	//display.Write(50,32)
	//
	//display.WriteBytes(130, []byte{1, 0, 1, 2})

	display.Write(1, 1, chip8.Sprite7)
	display.Write(50, 5, chip8.SpriteA)
	display.Write(25, 10, chip8.Sprite8)

	display.Run()
}
