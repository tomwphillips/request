package requester

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"cloud.google.com/go/storage"
)

func TestDecodeInstruction(t *testing.T) {
	in := []byte(`{"url": "http://google.com"}`)
	want := instruction{URL: "http://google.com"}
	got := decodeInstruction(in)
	if got != want {
		t.Errorf("decodeInstruction(%s) = %+v, want %+v", in, got, want)
	}
}

func TestGetURL(t *testing.T) {
	want := []byte("Hello, world!")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", want)
	}))
	in := instruction{URL: ts.URL}
	defer ts.Close()

	got := getURL(&in.URL)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("executeInstuction(%+v) = %s, want %s", in, got, want)
	}
}

func TestUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping upload test in short mode")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Creating client: %v", err)
	}

	bucketName := os.Getenv("REQUESTER_TEST_BUCKET")
	if bucketName == "" {
		t.Fatalf("REQUESTER_TEST_BUCKET environment variable not set")
	}

	objName := "test-object.txt"
	want := []byte("test bytes")
	err = upload(ctx, client, &want, objName, bucketName)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	r, err := client.Bucket(bucketName).Object(objName).NewReader(ctx)
	if err != nil {
		t.Fatalf("Failed reading object: %v", err)
	}
	defer r.Close()
	got, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Error reading from bucket: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %s, want %s", got, want)
	}
}
