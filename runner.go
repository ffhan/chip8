package chip8

import (
	"io"
)

func Run(reader io.Reader, version Version, display Display, keyboard Keyboard, speaker Speaker) {
	clock := NewClock()

	delayTimer := NewDelayTimer(clock)
	soundTimer := NewSoundTimer(speaker, clock)

	cpu := NewCPU(display, speaker, keyboard, clock, delayTimer, soundTimer, 12345)
	cpu.SetVersion(version)

	if err := cpu.LoadRom(reader); err != nil {
		panic(err)
	}

	go clock.Run(ClockFrequency)
	go cpu.Run()

	if err := display.Run(); err != nil {
		panic(err)
	}
}
