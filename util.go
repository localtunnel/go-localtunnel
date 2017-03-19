package localtunnel

import "sync"

// Atomic counter
type counter struct {
	m       sync.Mutex
	c       sync.Cond
	counter int
}

// Add value to counter
func (c *counter) Add(value int) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.c.L == nil {
		c.c.L = &c.m
	}

	c.counter += value
	c.c.Broadcast()
}

// WaitFor counter to reach target
func (c *counter) WaitFor(target int) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.c.L == nil {
		c.c.L = &c.m
	}

	for c.counter < target {
		c.c.Wait()
	}
}
