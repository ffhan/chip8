package chip8

import (
	"time"
)

const (
	ClockFrequency = 500 // Hz
)

type Clock struct {
	subs      []chan bool
	ticker    *time.Ticker
	frequency int64
}

func NewClock(frequency int64) *Clock {
	return &Clock{
		subs:      make([]chan bool, 0),
		frequency: frequency,
	}
}

func (c *Clock) Run() {
	c.ticker = time.NewTicker(time.Second / time.Duration(c.frequency))
	for range c.ticker.C {
		c.dispatch()
	}
}

func (c *Clock) Stop() {
	c.ticker.Stop()
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

func (c *Clock) Frequency() int64 {
	return c.frequency
}
