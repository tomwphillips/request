package bigquery

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
)

// GCSEvent is the payload of a Google Cloud Storage event
type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

const fileCreated string = "1"

// Read file from bucket
func read(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	if err := client.Close(); err != nil {
		return nil, err
	}
	r, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading %s/%s: %v", bucketName, objectName, err)
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

// ConsumeRequestOutput consumes a GCS event triggered the request cloud
// function writing to a GCS bucket.
func ConsumeRequestOutput(ctx context.Context, e GCSEvent) error {
	if e.ResourceState == "not_exists" {
		log.Printf("%v deleted", e.Name)
		return nil
	}
	if e.Metageneration == fileCreated {
		log.Printf("%v created", e.Name)
		return nil
	}
	log.Printf("%v metadata updated", e.Name)
	return nil
}
