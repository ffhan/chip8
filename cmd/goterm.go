package main

import (
	"chip8"
	"fmt"
)

func main() {
	display := chip8.NewDefaultGuiDisplay()
	//display.Write(10,32)
	//display.Write(128,32)
	//display.Write(50,32)
	//
	//display.WriteBytes(130, []byte{1, 0, 1, 2})

	fmt.Println(display.Write(1, 2, chip8.Sprite7))
	fmt.Println(display.Write(60, 5, chip8.SpriteA))
	fmt.Println(display.Write(25, 30, chip8.Sprite8))

	display.Run()
}
