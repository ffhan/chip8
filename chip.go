package chip8

import (
	"io"
)

func Run(reader io.Reader, version Version) {
	keyboard := NewDefaultKeyboard()

	display := NewDefaultGuiDisplay(keyboard)

	clock := NewClock()

	speaker := NewPulseAudioSpeaker(DefaultSpeakerFrequency)

	delayTimer := NewDelayTimer(clock)
	soundTimer := NewSoundTimer(speaker, clock)

	cpu := NewCPU(display, speaker, keyboard, clock, delayTimer, soundTimer, 12345)
	cpu.SetVersion(version)

	if err := cpu.LoadRom(reader); err != nil {
		panic(err)
	}

	go clock.Run(ClockFrequency)
	go cpu.Run()

	display.Run()
}
