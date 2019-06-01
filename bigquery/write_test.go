package bigquery

import (
	"bytes"
	"context"
	"log"
	"math"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
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
	if testing.Short() {
		t.Skip("Skipping StreamRecords test in short mode")
	}

	type Pet struct {
		ID     string  // needed because BiqQuery de-duplicates streaming inserts
		Name   string
		Animal string
		Age    int
	}

	records := []Pet{
		{uuid.New().String(), "Ronald", "Cat", 8},
		{uuid.New().String(), "Bessie", "Dog", 2},
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	defer client.Dataset(bqDataset).DeleteWithContents(ctx)
	if err != nil {
		t.Fatalf("creating client for %v: %v", projectID, err)
	}

	ds, err := InitializeDataset(ctx, client, bqDataset)
	if err != nil {
		t.Fatalf("initializing dataset: %v", err)
	}

	th, err := InitializeTable(ctx, ds, bqTable, Pet{})
	if err != nil {
		t.Fatalf("creating table: %v", err)
	}

	err = StreamRecords(ctx, th, records)
	if err != nil {
		t.Fatalf("StreamRecord != nil, got %v", err)
	}

	// Streaming inserts go into a buffer first, not the table directly, but take a while to appear
	var md *bigquery.TableMetadata
	retries := 5
	for i := 1; i <= retries; i++ {
		time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))

		if md, err = th.Metadata(ctx); err != nil {
			t.Fatalf("metadata error: %v", err)
		}

		if i <= retries && md.StreamingBuffer == nil {
			continue  // no buffer yet, retry
		} else if md.StreamingBuffer == nil {
			t.Fatalf("no streaming buffer")
		}

		got := int(md.StreamingBuffer.EstimatedRows)
		if i <= retries && got != len(records) {
			continue  // not enough records yet, retry
		} else if got != len(records) {
			t.Fatalf("Streaming buffer = %v, want %v", got, len(records))
		}
	}
}
