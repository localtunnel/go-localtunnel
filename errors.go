package localtunnel

import "errors"

// ErrListenerClosed indicates that the listener as closed
var ErrListenerClosed = errors.New("listener was closed")
