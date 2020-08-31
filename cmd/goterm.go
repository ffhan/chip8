package main

import (
	"chip8"
	"fmt"
	"time"
)

func main() {
	display := chip8.NewDefaultGuiDisplay()
	fmt.Println(display.Write(1, 2, chip8.Sprite7))
	fmt.Println(display.Write(60, 5, chip8.SpriteA))
	fmt.Println(display.Write(25, 30, chip8.Sprite8))

	go func() {
		speaker := chip8.NewPulseAudioSpeaker(chip8.DefaultSpeakerFrequency)
		speaker.Play()

		time.Sleep(2 * time.Second)
		speaker.Stop()
		fmt.Println("stopped")
		time.Sleep(2 * time.Second)
		speaker.Play()
		fmt.Println("started")

		time.Sleep(2 * time.Second)

		speaker.Stop()
		fmt.Println("stopped")

		time.Sleep(5 * time.Second)
		fmt.Println("exiting")
	}()

	display.Run()
}
