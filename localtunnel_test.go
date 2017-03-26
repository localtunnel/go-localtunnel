package localtunnel

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type testLog struct {
	*testing.T
}

func (t testLog) Println(v ...interface{}) {
	t.Log(v...)
}

// testLocalTunnel against server at baseURL (if given)
func testLocalTunnel(baseURL string, t *testing.T) {
	// Collect some random data
	random := make([]byte, 4*1024*1024) // 4 MiB
	if _, err := rand.Read(random); err != nil {
		t.Fatal("failed to generate random data, error:", err)
	}

	// Setup a test server that we wish to expose
	port := 60000
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if r.URL.Query().Get("random") == "true" {
				w.Write(random)
			} else {
				w.Write([]byte("hello-world"))
			}
		}),
	}
	go server.ListenAndServe()
	defer server.Close()

	t.Log("setting up LocalTunnel")
	lt, err := New(port, "", Options{
		Log:            testLog{t},
		MaxConnections: 2,
		BaseURL:        baseURL,
	})
	if err != nil {
		t.Fatal("failed to create LocalTunnel, error: ", err)
	}

	// Sleep for 3s giving the server time to register
	time.Sleep(3 * time.Second)

	// Make http.Client with timeout of 30 seconds
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Let's make 3 requests for good measure
	for i := 0; i < 3; i++ {
		t.Log("sending test request to:", lt.URL())
		var res *http.Response
		res, err = client.Get(lt.URL())
		if err != nil {
			t.Fatal("failed to send GET request through tunnel, error: ", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Error("expected 200 ok, got status: ", res.StatusCode)
		}

		var data []byte
		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal("failed to read response from tunnel, error: ", err)
		}
		if string(data) != "hello-world" {
			t.Error("unexpected response, data: ", string(data))
		}
	}

	t.Log("sending testing request to:", lt.URL()+"/?random=true")
	var res *http.Response
	res, err = client.Get(lt.URL() + "/?random=true")
	if err != nil {
		t.Fatal("failed to send GET request through tunnel, error: ", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Error("expected 200 ok, got status: ", res.StatusCode)
	}

	var data []byte
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("failed to read response from tunnel, error: ", err)
	}
	if bytes.Compare(data, random) != 0 {
		t.Error("unexpected random response, size:", len(data))
	}

	t.Log("closing LocalTunnel")
	err = lt.Close()
	if err != nil {
		t.Error("error closing the tunnel: ", err)
	}
}

func TestLocalTunnel(t *testing.T) {
	testLocalTunnel("", t)
}
