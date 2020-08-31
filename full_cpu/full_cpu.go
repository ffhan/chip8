package main

import "chip8"

func main() {
	display := chip8.NewDefaultGuiDisplay()
	cpu := chip8.NewCPU(display, chip8.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency), chip8.NewFyneKeyboard(display.GetCanvas()), 12345)

	cpu.Run()
}
