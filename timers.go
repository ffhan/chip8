package chip8

// FIXME: timers won't stop when the whole execution has to stop - rethink tickers (use clean goroutines)

const (
	delayFrequency = 60 // Hz
)

type DelayTimer struct {
	dt byte // delay ticker

	sub <-chan bool
}

func NewDelayTimer(clock *Clock) *DelayTimer {
	sub := clock.Subscribe()
	d := &DelayTimer{sub: sub}
	go d.serviceDelay()
	return d
}

func (d *DelayTimer) Get() byte {
	return d.dt
}

func (d *DelayTimer) Set(delay byte) {
	d.dt = delay
}

func (d *DelayTimer) serviceDelay() {
	for {
		<-d.sub
		if d.dt == 0 {
			continue
		}
		d.dt -= 1
	}
}

type SoundTimer struct {
	st byte // sound ticker

	speaker Speaker
	sub     <-chan bool
}

func NewSoundTimer(speaker Speaker, clock *Clock) *SoundTimer {
	sub := clock.Subscribe()
	s := &SoundTimer{speaker: speaker, sub: sub}
	go s.serviceDelay()
	return s
}

func (s *SoundTimer) Get() byte {
	return s.st
}

func (s *SoundTimer) Set(delay byte) {
	s.st = delay
}

func (s *SoundTimer) serviceDelay() {
	playing := false
	for {
		<-s.sub
		if s.st == 0 {
			if playing {
				s.speaker.Stop()
				playing = false
			}
			continue
		}
		s.st -= 1
		if !playing {
			s.speaker.Play()
			playing = true
		}
	}
}
