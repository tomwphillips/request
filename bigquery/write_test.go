package bigquery

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
)

var projectID = os.Getenv("REQUEST_GCP_TEST_PROJECT")
var bqDataset = os.Getenv("REQUEST_BQ_TEST_DATASET")
var bqTable = os.Getenv("REQUEST_BQ_TEST_TABLE")

func captureOutput(f func()) string {
	var buf bytes.Buffer
	originalFlags := log.Flags()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	log.SetFlags(originalFlags)
	return buf.String()
}

func TestGCSWriteEvent(t *testing.T) {
	ctx := context.Background()
	var tests = []struct {
		e    GCSEvent
		log  string
		want bool
	}{
		{
			GCSEvent{"", "filename", "", "not_exists"},
			"filename deleted\n",
			false,
		},
		{
			GCSEvent{"", "filename", fileCreated, ""},
			"filename created\n",
			true,
		},
		{
			GCSEvent{"", "filename", "", ""},
			"filename metadata updated\n",
			false,
		},
	}

	for _, test := range tests {
		var got bool
		gotLog := captureOutput(func() { got = GCSWriteEvent(ctx, test.e) })
		if test.log != gotLog {
			t.Errorf("GCSWriteEvent(ctx, %v) logged %+v, want %+v", test.e, gotLog, test.log)
		}
		if got != test.want {
			t.Errorf("GCSWriteEvent(ctx, %v) = %+v, want %+v", test.e, got, test.want)
		}
	}
}

func TestStreamRecords(t *testing.T) {
	type Pet struct {
		Name   string
		Animal string
		Age    int
	}
	records := []Pet{
		{"Ronald", "Cat", 8},
		{"Bessie", "Dog", 2},
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		t.Errorf("creating client for %v: %v", projectID, err)
	}

	ds, err := InitializeDataset(ctx, client, bqDataset)
	if err != nil {
		t.Errorf("initializing dataset: %v", err)
	}

	th, err := InitializeTable(ctx, ds, bqTable, Pet{})
	if err != nil {
		t.Errorf("creating table: %v", err)
	}

	err = StreamRecords(ctx, th, records)
	if err != nil {
		t.Errorf("StreamRecord != nil, got %v", err)
	}

	// don't do a count on table because streaming inserts go into a buffer, not the table
	md, _ := th.Metadata(ctx)
	got := int(md.StreamingBuffer.EstimatedRows)
	if got != len(records) {
		t.Errorf("streaming buffer estimate = %v, not %v", got, len(records))
	}

	// TODO: defer this
	if err := client.Dataset(bqDataset).DeleteWithContents(ctx); err != nil {
		t.Errorf("Deleting test dataset: %v", err)
	}
}
