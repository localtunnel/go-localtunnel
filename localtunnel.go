package localtunnel

import (
	"fmt"
	"io"
	"net"
)

// LocalTunnel forwards remote requests to a port on localhost
type LocalTunnel struct {
	listener  *Listener
	localAddr string
}

// New returns a LocalTunnel forwarding requests to port on host
//
// host defaults to 'localhost', and options defaults to using localtunnel.me
func New(port int, host string, options Options) (*LocalTunnel, error) {
	if host == "" {
		host = "localhost"
	}

	l, err := Listen(options)
	if err != nil {
		return nil, err
	}

	lt := &LocalTunnel{
		listener:  l,
		localAddr: fmt.Sprintf("%s:%d", host, port),
	}
	go lt.listen()
	return lt, nil
}

// URL returns the URL at which the localtunnel is exposed
func (lt *LocalTunnel) URL() string {
	return lt.listener.URL()
}

func (lt *LocalTunnel) listen() {
	for {
		remoteConn, err := lt.listener.Accept()
		if err != nil {
			break
		}

		go lt.forward(remoteConn)
	}
}

func (lt *LocalTunnel) forward(remoteConn net.Conn) {
	localConn, err := net.Dial("tcp", lt.localAddr)
	if err != nil {
		remoteConn.Close()
		return
	}

	go func() {
		io.Copy(remoteConn, localConn)
		remoteConn.Close()
	}()
	go func() {
		io.Copy(localConn, remoteConn)
		localConn.Close()
	}()
}

// Close the localtunnel aborting all connections
func (lt *LocalTunnel) Close() error {
	return lt.listener.Close()
}
