package chip8

import (
	"io"
)

func Run(reader io.Reader, version Version, display Display, keyboard Keyboard, speaker Speaker) func() {
	clockFreq := ClockFrequency
	if version != Chip8 {
		clockFreq *= 2
	}

	clock := NewClock(int64(clockFreq))

	delayTimer := NewDelayTimer(clock)
	soundTimer := NewSoundTimer(speaker, clock)

	cpu := NewCPU(display, speaker, keyboard, clock, clock, delayTimer, soundTimer, nil)
	cpu.SetVersion(version)

	if err := cpu.LoadRom(reader); err != nil {
		panic(err)
	}

	go clock.Run()
	go cpu.Run()

	if err := display.Run(); err != nil {
		panic(err)
	}

	return func() {
		cpu.Stop()
	}
}
