package main

import (
	"chip8"
	"chip8/fyne"
	"chip8/pulseaudio"
	"flag"
	"os"
)

func main() {
	romPath := flag.Arg(0)
	if romPath == "" {
		romPath = "BC_test.ch8"
	}

	keyboard := chip8.NewDefaultKeyboard()

	display := fyne.NewDefaultGuiDisplay(keyboard)
	speaker := pulseaudio.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency)

	rom, err := os.Open(romPath)
	if err != nil {
		panic(err)
	}
	defer rom.Close()

	chip8.Run(rom, chip8.SuperChip48, display, keyboard, speaker)
}
