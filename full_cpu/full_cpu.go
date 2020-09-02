package main

import (
	"chip8"
	"chip8/fyne"
	"chip8/pulseaudio"
	"os"
)

func main() {
	keyboard := chip8.NewDefaultKeyboard()

	display := fyne.NewDefaultGuiDisplay(keyboard)
	speaker := pulseaudio.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency)

	rom, err := os.Open("chip8-/roms/Space Invaders [David Winter].ch8")
	if err != nil {
		panic(err)
	}
	defer rom.Close()

	chip8.Run(rom, chip8.SuperChip48, display, keyboard, speaker)
}
