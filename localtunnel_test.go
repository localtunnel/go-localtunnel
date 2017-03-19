package localtunnel

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type testLog struct {
	*testing.T
}

func (t testLog) Println(v ...interface{}) {
	t.Log(v...)
}

func TestLocalTunnel(t *testing.T) {
	port := 60000
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello-world"))
		}),
	}
	go server.ListenAndServe()
	defer server.Close()

	t.Log("setting up LocalTunnel")
	lt, err := New(port, "", Options{Log: testLog{t}, MaxConnections: 1})
	if err != nil {
		t.Fatal("failed to create LocalTunnel, error: ", err)
	}

	t.Log("sending test request")
	res, err := http.Get(lt.URL())
	if err != nil {
		t.Fatal("failed to send GET request through tunnel, error: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("expected 200 ok, got status: ", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("failed to read response from tunnel, error: ", err)
	}
	if string(data) != "hello-world" {
		t.Error("unexpected response, data: ", string(data))
	}

	t.Log("closing LocalTunnel")
	err = lt.Close()
	if err != nil {
		t.Error("error closing the tunnel: ", err)
	}
}
