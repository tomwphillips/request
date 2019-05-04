package request

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"cloud.google.com/go/storage"
)

var bucketName = os.Getenv("REQUESTER_TEST_BUCKET")

func init() {
	if bucketName == "" {
		log.Fatalf("REQUESTER_TEST_BUCKET environment variable not set")
	}
}

func TestDecodeInstruction(t *testing.T) {
	in := []byte(`{"url": "http://google.com", "bucket": "bucket-name"}`)
	want := instruction{URL: "http://google.com", Bucket: "bucket-name"}
	got, _ := decodeInstruction(in)
	if got != want {
		t.Errorf("decodeInstruction(%s) = %+v, want %+v", in, got, want)
	}
}

func TestGetURL(t *testing.T) {
	want := []byte("Hello, world!")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", want)
	}))
	defer ts.Close()

	in := ts.URL
	got, _ := getURL(in)
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

	objName := "test-object.txt"
	want := []byte("test bytes")
	obj, err := upload(ctx, client, &want, objName, bucketName)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		t.Fatalf("Reader failed: %v", err)
	}
	defer r.Close()
	got, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got %s, want %s", got, want)
	}
}

func TestExecute(t *testing.T) {
	want := []byte("Hello, world!")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", want)
	}))
	defer ts.Close()

	ctx := context.Background()
	in := instruction{URL: ts.URL, Bucket: bucketName}
	obj, err := execute(ctx, in)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	r, err := obj.NewReader(ctx)
	defer r.Close()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	got, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Read from object: %s, want %s", got, want)
	}
}
