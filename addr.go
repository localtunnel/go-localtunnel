package localtunnel

// Network return "localtunnel"
func (a Addr) Network() string {
	return "localtunnel"
}

// String returns the URL
func (a Addr) String() string {
	return a.URL
}

// URL returns the URL that the listener is exposed on.
func (l *Listener) URL() string {
	return l.url
}
