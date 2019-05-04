package bigquery

import (
	"log"
	"context"
)

// GCSEvent is the payload of a Google Cloud Storage event
type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

const fileCreated string = "1"

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
