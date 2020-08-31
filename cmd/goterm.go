package main

import (
	"chip8"
	"fmt"
	"time"
)

func main() {
	keyboard := chip8.NewDefaultKeyboard()
	display := chip8.NewDefaultGuiDisplay(keyboard)
	fmt.Println(display.Write(1, 2, chip8.Sprite7))
	fmt.Println(display.Write(60, 5, chip8.SpriteA))
	fmt.Println(display.Write(25, 30, chip8.Sprite8))

	go func() {
		speaker := chip8.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency)
		speaker.Play()

		time.Sleep(2 * time.Second)
		speaker.Stop()
		time.Sleep(2 * time.Second)
		speaker.Play()
		time.Sleep(2 * time.Second)
		speaker.Stop()
		time.Sleep(5 * time.Second)
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("is A down?", keyboard.IsDown(chip8.KeyA))
			fmt.Println("got key ", keyboard.Wait())
		}
	}()

	display.Run()
}
