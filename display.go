package chip8

type Display interface {
	Clear()
	Write(x, y byte, bytes []byte) bool
	Run() error
	Repaint()
}
