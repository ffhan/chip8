package chip8

import "testing"

type fakeSpeaker struct {
}

func (f fakeSpeaker) Play(freq float32) {
	panic("implement me")
}

func (f fakeSpeaker) Stop() {
	panic("implement me")
}

type fakeKeyboard struct {
}

func (f fakeKeyboard) Wait() Key {
	panic("implement me")
}

func (f fakeKeyboard) IsDown(key Key) bool {
	panic("implement me")
}

func TestCPU_Step(t *testing.T) {
	display := NewDefaultGuiDisplay()
	display.updateScreen()
}
