# Request

`request` is a package containing functions get the contents of a URL and write it to [Google Cloud Storage](https://cloud.google.com/storage/) (GCS).

The `request/bigquery` package contains functions to read the files written to GCS and stream them to [BigQuery](https://cloud.google.com/bigquery/).

Both of these are intended for use as serverless functions, e.g. Google [Cloud Functions](https://cloud.google.com/functions/). The purpose is to regularly collect data from a URL and store it for further analysis. See [`tfl`](http://github.com/tomwphillips/tfl) for an example.

This is my first Go project, so it might not be entirely idiomatic.

## Tests

Set the following environment variables (with appropriate values):

```
REQUEST_TEST_BUCKET="name-of-test-bucket"
REQUEST_GCP_TEST_PROJECT="name-of-project"
REQUEST_BQ_TEST_DATASET="request_test"
REQUEST_BQ_TEST_TABLE="test"
```

Then:

```
go test ./...
```

Some of them are slow integration tests. You can skip them with the `-short` flag.
