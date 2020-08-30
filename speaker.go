package chip8

const (
	DefaultSpeakerFrequency = 293.6647
)

type Speaker interface {
	Play(freq float32)
	Stop()
}
