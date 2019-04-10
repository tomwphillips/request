package requester

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeInstruction(t *testing.T) {
	in := []byte(`{"url": "http://google.com"}`)
	want := instruction{URL: "http://google.com"}
	got := decodeInstruction(in)
	if got != want {
		t.Errorf("decodeInstruction(%s) = %+v, want %+v", in, got, want)
	}
}

func TestExecuteInstruction(t *testing.T) {
	want := "Hello, world!"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	in := instruction{URL: ts.URL}
	defer ts.Close()

	got := executeInstruction(in)
	if want != got {
		t.Errorf("executeInstuction(%+v) = %s, want %s", in, got, want)
	}
}
