package chip8

import (
	"io"
)

func Run(reader io.Reader, version Version, display Display, keyboard Keyboard, speaker Speaker) {
	clock := NewClock()
	timerClock := NewClock()

	delayTimer := NewDelayTimer(timerClock)
	soundTimer := NewSoundTimer(speaker, timerClock)

	cpu := NewCPU(display, speaker, keyboard, clock, timerClock, delayTimer, soundTimer, nil)
	cpu.SetVersion(version)

	if err := cpu.LoadRom(reader); err != nil {
		panic(err)
	}

	clockFreq := ClockFrequency
	timerFreq := TimerFrequency
	if version != Chip8 {
		clockFreq *= 2
	}

	go clock.Run(int64(clockFreq))
	go timerClock.Run(int64(timerFreq))
	go cpu.Run()

	if err := display.Run(); err != nil {
		panic(err)
	}
}
