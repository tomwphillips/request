package bigquery

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/bigquery"
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

// GCSWriteEvent returns true if event describes creation of a new file.
func GCSWriteEvent(ctx context.Context, e GCSEvent) bool {
	if e.ResourceState == "not_exists" {
		log.Printf("%v deleted", e.Name)
		return false
	}
	if e.Metageneration == fileCreated {
		log.Printf("%v created", e.Name)
		return true
	}
	log.Printf("%v metadata updated", e.Name)
	return false
}

// InitializeDataset returns handle to dataset. Creates dataset if it doesn't exist.
func InitializeDataset(ctx context.Context, client *bigquery.Client, id string) (*bigquery.Dataset, error) {
	ds := client.Dataset(id)
	if _, err := ds.Metadata(ctx); err != nil {
		if err := ds.Create(ctx, &bigquery.DatasetMetadata{Location: "EU"}); err != nil {
			return nil, fmt.Errorf("creating dataset %v: %v", id, err)
		}
	}
	return ds, nil
}

// InitializeTable returns handle to table. Checks record matches table schema.
// Creates table if it doesn't exist.
func InitializeTable(ctx context.Context, ds *bigquery.Dataset, table string, record interface{}) (*bigquery.Table, error) {
	t := ds.Table(table)
	if _, err := t.Metadata(ctx); err == nil {
		// TODO: table already exists, check schema matches otherwise return error
		return t, nil
	}
	schema, err := bigquery.InferSchema(record)
	if err != nil {
		return nil, fmt.Errorf("infering schema from %+v: %v", record, err)
	}

	if err := t.Create(ctx,
		&bigquery.TableMetadata{
			Name:   table,
			Schema: schema,
		}); err != nil {
		return nil, fmt.Errorf("creating table %v: %v", table, err)
	}
	return ds.Table(table), nil
}

// StreamRecords to BigQuery.
func StreamRecords(ctx context.Context, t *bigquery.Table, records interface{}) error {
	ins := t.Inserter()
	return ins.Put(ctx, records)
}
