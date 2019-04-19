package requester

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

type instruction struct {
	URL string
}

// Upload writes bytes to an object in a Google Storage bucket
func upload(ctx context.Context, client *storage.Client, b *[]byte, object string, bucket string) error {
	bh := *client.Bucket(bucket)

	if _, err := bh.Attrs(ctx); err != nil {
		return err // bucket doesn't exist
	}

	obj := bh.Object(object)
	w := obj.NewWriter(ctx)
	if _, err := w.Write(*b); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func decodeInstruction(m []byte) instruction {
	var i instruction
	err := json.Unmarshal(m, &i)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}
	return i
}

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
}

func executeInstruction(i instruction) string {
	return getURL(&i.URL)
}
