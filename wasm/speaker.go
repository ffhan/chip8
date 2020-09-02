package wasm

import "syscall/js"

type speaker struct {
	playMusic js.Value
	stopMusic js.Value
}

func NewSpeaker() *speaker {
	return &speaker{
		playMusic: js.Global().Get("playMusic"),
		stopMusic: js.Global().Get("stopMusic"),
	}
}

func (s *speaker) Play() {
	s.playMusic.Invoke()
}

func (s *speaker) Stop() {
	s.stopMusic.Invoke()
}
