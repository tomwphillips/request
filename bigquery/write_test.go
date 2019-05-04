package bigquery

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"
)

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

func TestConsumeRequestOutput(t *testing.T) {
	ctx := context.Background()
	var tests = []struct{
		e GCSEvent
		want string
	}{
		{
			GCSEvent{"", "filename", "", "not_exists"},
			"filename deleted\n",
		},
		{
			GCSEvent{"", "filename", fileCreated, ""},
			"filename created\n",
		},
		{
			GCSEvent{"", "filename", "", ""},
			"filename metadata updated\n",
		},
	}

	for _, test := range tests {
		got := captureOutput(func() {ConsumeRequestOutput(ctx, test.e)})
		if test.want != got {
			t.Errorf("ConsumeRequestOutput(ctx, %v) = %+v, want %+v", test.e, got, test.want)
		}
	}
}
