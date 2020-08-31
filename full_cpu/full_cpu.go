package main

import (
	"chip8"
	"os"
)

func main() {
	keyboard := chip8.NewDefaultKeyboard()
	display := chip8.NewDefaultGuiDisplay(keyboard)
	cpu := chip8.NewCPU(display, chip8.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency), keyboard, 12345)

	open, err := os.Open("c8games/CONNECT4")
	if err != nil {
		panic(err)
	}

	if err = cpu.LoadRom(open); err != nil {
		panic(err)
	}

	cpu.Run()
}
