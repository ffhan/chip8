package chip8

import (
	"fmt"
	"sync"
	"time"
)

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
	KeyDown(r rune)
	KeyUp(r rune)
}

type defaultKeyboard struct {
	pressed sync.Map
	events  chan Key
}

func (d *defaultKeyboard) MapRuneToKey(r rune) Key {
	switch r {
	case 'A':
		return KeyA
	case 'B':
		return KeyB
	case 'C':
		return KeyC
	case 'D':
		return KeyD
	case 'E':
		return KeyE
	case 'F':
		return KeyF
	case '0':
		return Key0
	case '1':
		return Key1
	case '2':
		return Key2
	case '3':
		return Key3
	case '4':
		return Key4
	case '5':
		return Key5
	case '6':
		return Key6
	case '7':
		return Key7
	case '8':
		return Key8
	case '9':
		return Key9
	}
	return Key0
}

func (d *defaultKeyboard) KeyDown(r rune) {
	fmt.Println("down ", string(r))
	key := d.MapRuneToKey(r)
	d.pressed.Store(key, true)
	select {
	case d.events <- key:
		return
	case <-time.After(10 * time.Microsecond):
		return
	}
}

func (d *defaultKeyboard) KeyUp(r rune) {
	d.pressed.Store(d.MapRuneToKey(r), false)
}

func NewDefaultKeyboard() *defaultKeyboard {
	return &defaultKeyboard{
		pressed: sync.Map{},
		events:  make(chan Key),
	}
}

func (d *defaultKeyboard) Wait() Key {
	e := <-d.events
	return e
}

func (d *defaultKeyboard) IsDown(key Key) bool {
	val, _ := d.pressed.LoadOrStore(key, false)
	return val.(bool)
}
