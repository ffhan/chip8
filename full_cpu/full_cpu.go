package main

import (
	"chip8"
	"os"
)

func main() {
	keyboard := chip8.NewDefaultKeyboard()
	display := chip8.NewDefaultGuiDisplay(keyboard)
	cpu := chip8.NewCPU(display, chip8.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency), keyboard, 12345)

	open, err := os.Open("c8games/TETRIS")
	if err != nil {
		panic(err)
	}

	if err = cpu.LoadRom(open); err != nil {
		panic(err)
	}

	go cpu.Run()

	display.Run()
}
