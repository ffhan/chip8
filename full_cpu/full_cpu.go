package main

import (
	"chip8"
	"os"
)

func main() {
	rom, err := os.Open("chip8-/roms/Space Invaders [David Winter].ch8")
	if err != nil {
		panic(err)
	}
	defer rom.Close()

	chip8.Run(rom, chip8.SuperChip48)
}
