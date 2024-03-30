package util

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	oe := errors.New("graceful exit")
	e2 := fmt.Errorf("%w err", ErrGracefulExit)
	t.Log(errors.Is(e2, ErrGracefulExit), errors.Is(e2, oe))
}

func TestClient(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	// io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	<-time.After(time.Second)
}
