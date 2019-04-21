package requester

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

type instruction struct {
	URL    string
	Bucket string
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// decodeInstruction from JSON-encoded byte array
func decodeInstruction(m []byte) instruction {
	var i instruction
	err := json.Unmarshal(m, &i)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}
	return i
}

// getURL returns contents at URL
func getURL(url *string) []byte {
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalf("Get request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read response body failed: %v", err)
	}
	return body
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
		return nil, err // bucket doesn't exist
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

// ConsumePubSub decodes and execute instructions
func ConsumePubSub(ctx context.Context, m PubSubMessage) (*storage.ObjectHandle, error) {
	i := decodeInstruction(m.Data)
	body := getURL(&i.URL)
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	object := fmt.Sprintf("%x", getHash(&body))
	return upload(ctx, client, &body, object, i.Bucket)
}
