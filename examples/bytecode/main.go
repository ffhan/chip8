package main

import (
	"bytes"
	"chip8"
	"chip8/fyne"
	"chip8/pulseaudio"
	"strings"
)

func main() {
	keyboard := chip8.NewDefaultKeyboard()

	display := fyne.NewDefaultGuiDisplay(keyboard)
	speaker := pulseaudio.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency)

	reader := strings.NewReader("6010 6110 FA29 D015 00FD")
	codes, err := chip8.ParseCodes(reader)
	if err != nil {
		panic(err)
	}
	chip8.Run(bytes.NewBuffer(codes), chip8.SuperChip48, display, keyboard, speaker)
}
