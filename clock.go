package chip8

import (
	"time"
)

const (
	clockFrequency = 60 // Hz
)

type Clock struct {
	subs     []chan bool
	lastStep time.Time
}

func (c *Clock) Step() {
	now := time.Now()
	if now.Sub(c.lastStep) >= (time.Second / clockFrequency) {
		c.dispatch()
		c.lastStep = now
	}
}

func (c *Clock) dispatch() {
	for _, sub := range c.subs {
		select {
		case sub <- true:
			continue
		case <-time.After(20 * time.Microsecond):
			continue
		}
	}
}

func (c *Clock) Subscribe() <-chan bool {
	sub := make(chan bool)
	c.subs = append(c.subs, sub)
	return sub
}
