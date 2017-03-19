LocalTunnel Client Library for Go
=================================

A [localtunnel.me](https://localtunnel.me) client library exposing localtunnel
connections through a `net.Listener` implementation. While localtunnel only
supports forwarding HTTP(S) connections, this is useful as you can hook it up
to `http.Server` directly. Neat, if writing test-suites or command-line
utilities exposing web-hooks of localtunnel.

```go
// Setup a listener for localtunnel
listener, err := localtunnel.Listen(localtunnel.Options{})

// Create your server...
server := http.Server{
    Handler: http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        ...
    })
}

// Handle request from localtunnel
server.Serve(listener)
```

See [documentation](https://godoc.org/github.com/jonasfj/go-localtunnel) for
more details.


License
-------
This package is released under [MPLv2](https://www.mozilla.org/en-US/MPL/2.0/).
