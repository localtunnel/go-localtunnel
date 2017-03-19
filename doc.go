// Package localtunnel implements a client library for https://localtunnel.me
//
// In addition to providing the LocalTunnel client which will forward requests
// from subdomain.localtunnel.me to a port on localhost. This package also
// provides an implementation of net.Listener which exposes connections from
// localtunnel. This enables users to serve http requests directly, without
// listening to a port on localhost.
//
//   // Setup a listener for localtunnel
//   listener, err := localtunnel.Listen(localtunnel.Options{})
//
//   // Create your server...
//   server := http.Server{
//       Handler: http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
//           w.WriteHeader(200)
//           ...
//       })
//   }
//
//   // Handle request from localtunnel
//   server.Serve(listener)
package localtunnel
