package chip8

const (
	TimerFrequency = 60 // Hz
)

type DelayTimer struct {
	dt byte // delay ticker

	sub   <-chan bool
	clock *Clock
}

func NewDelayTimer(clock *Clock) *DelayTimer {
	sub := clock.Subscribe()
	d := &DelayTimer{sub: sub, clock: clock}
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
	steps := d.clock.Frequency() / TimerFrequency
	for {
		for i := int64(0); i < steps; i++ {
			<-d.sub
		}
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
	clock   *Clock
}

func NewSoundTimer(speaker Speaker, clock *Clock) *SoundTimer {
	sub := clock.Subscribe()
	s := &SoundTimer{speaker: speaker, sub: sub, clock: clock}
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
	steps := s.clock.Frequency() / TimerFrequency
	for {
		for i := int64(0); i < steps; i++ {
			<-s.sub
		}
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
