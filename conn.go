package localtunnel

import (
	"net"
	"sync"
)

type conn struct {
	net.Conn
	Buffer   [1]byte
	Done     chan<- struct{}
	onceDone sync.Once
	mRead    sync.Mutex
	read     bool
}

func (c *conn) done() {
	c.onceDone.Do(func() {
		close(c.Done)
	})
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.mRead.Lock()
	defer c.mRead.Unlock()

	if c.read {
		n, err = c.Conn.Read(b)
		if err != nil {
			c.done()
		}
		return
	}

	if len(b) == 0 {
		return 0, nil
	}
	c.read = true
	b[0] = c.Buffer[0]
	if len(b) > 1 {
		n, err = c.Conn.Read(b[1:])
		if err != nil {
			c.done()
		}
	}
	n++
	return
}

func (c *conn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		c.done()
	}
	return
}

func (c *conn) Close() error {
	defer c.done()
	return c.Conn.Close()
}
