package chip8

type Key uint8

const (
	Key0 = iota
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
)

type Keyboard interface {
	Wait() Key
	IsDown(key Key) bool
}
