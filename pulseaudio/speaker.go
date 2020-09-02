package pulseaudio

import "github.com/faiface/beep/speaker"

type sound struct {
}

func (s sound) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		if i%1024 == 0 {
			samples[i] = [2]float64{0.5, 0.5}
		}
	}
	return len(samples), true
}

func (s sound) Err() error {
	panic("implement me")
}

type pulseAudioSpeaker struct {
}

func NewPulseAudioSpeaker(freq float64) *pulseAudioSpeaker {
	if err := speaker.Init(44100, 4096); err != nil {
		panic(err)
	}
	return &pulseAudioSpeaker{}
}

func (p *pulseAudioSpeaker) Play() {
	speaker.Play(sound{})
}

func (p *pulseAudioSpeaker) Stop() {
	speaker.Clear()
}
