package chip8

import (
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

var (
	keyMappings = map[rune]Key{
		'a': KeyA, 'b': KeyB, 'c': KeyC, 'd': KeyD,
		'e': KeyE, 'f': KeyF, '0': Key0, '1': Key1,
		'2': Key2, '3': Key3, '4': Key4, '5': Key5,
		'6': Key6, '7': Key7, '8': Key8, '9': Key9,
	}
	reverseKeyMappings = map[Key]rune{
		KeyA: 'A', KeyB: 'B', KeyC: 'C', KeyD: 'D',
		KeyE: 'E', KeyF: 'F', Key0: '0', Key1: '1',
		Key2: '2', Key3: '3', Key4: '4', Key5: '5',
		Key6: '6', Key7: '7', Key8: '8', Key9: '9',
	}
)

type defaultKeyboard struct {
	pressed      sync.Map
	events       chan Key
	searchedKeys map[string]bool
}

func (d *defaultKeyboard) MapRuneToKey(r rune) Key {
	if val, ok := keyMappings[r]; ok {
		return val
	}
	return Key0
}

func (d *defaultKeyboard) MapKeyToRune(k Key) rune {
	if val, ok := reverseKeyMappings[k]; ok {
		return val
	}
	return '0'
}

func (d *defaultKeyboard) KeyDown(r rune) {
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
	d := &defaultKeyboard{
		pressed:      sync.Map{},
		events:       make(chan Key),
		searchedKeys: make(map[string]bool),
	}
	//go func() {
	//	t := time.NewTicker(1 * time.Second)
	//	for range t.C {
	//		fmt.Println(d.searchedKeys)
	//	}
	//}()
	return d
}

func (d *defaultKeyboard) Wait() Key {
	e := <-d.events
	return e
}

func (d *defaultKeyboard) IsDown(key Key) bool {
	d.searchedKeys[string(d.MapKeyToRune(key))] = true
	val, _ := d.pressed.LoadOrStore(key, false)
	return val.(bool)
}
