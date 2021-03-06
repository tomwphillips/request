package request

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/storage"
)

// Instruction to get URL and destination bucket
type Instruction struct {
	URL    string
	Bucket string
}

// DecodeInstruction from JSON-encoded byte array
func DecodeInstruction(m []byte) (Instruction, error) {
	var i Instruction
	err := json.Unmarshal(m, &i)
	return i, err
}

// GetURL returns contents at URL
func getURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getting %s: %v", url, err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getHash(b *[]byte) []byte {
	h := sha1.New()
	h.Write(*b)
	return h.Sum(nil)
}

// Upload writes bytes to an object in a Google Storage bucket
func upload(ctx context.Context, client *storage.Client, b *[]byte, object string, bucket string) (*storage.ObjectHandle, error) {
	bh := *client.Bucket(bucket)
	if _, err := bh.Attrs(ctx); err != nil {
		return nil, fmt.Errorf("getting bucket %s metadata: %v", bucket, err)
	}

	obj := bh.Object(object)
	w := obj.NewWriter(ctx)
	if _, err := w.Write(*b); err != nil {
		return obj, err
	}
	if err := w.Close(); err != nil {
		return obj, err
	}

	return obj, nil
}

// Execute instruction to get contents of a URL and upload to GCS
func Execute(ctx context.Context, i Instruction) (*storage.ObjectHandle, error) {
	body, err := getURL(i.URL)
	if err != nil {
		return nil, err
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	object := fmt.Sprintf("%x", getHash(&body))
	return upload(ctx, client, &body, object, i.Bucket)
}
