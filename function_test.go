package requester

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestConsumePubSub(t *testing.T) {
	in := `{"url": "http://google.com"}`
	want := "Received instruction: {URL:http://google.com}\n"

	r, w, _ := os.Pipe()
	log.SetOutput(w)
	originalFlags := log.Flags()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	m := PubSubMessage{
		Data: []byte(in),
	}
	ConsumePubSub(context.Background(), m)

	w.Close()
	log.SetOutput(os.Stderr)
	log.SetFlags(originalFlags)

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if got := string(out); got != want {
		t.Errorf("ConsumePubSub(%q) = %q, want %q", in, got, want)
	}
}
