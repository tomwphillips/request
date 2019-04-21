# Requester

A GCP Cloud Function triggered via PubSub that makes a HTTP GET request to a URL and saves the response body to a GCP Storage Bucket. Intended for getting data from APIs at regular intervals.

My first Go project, so it's probably not idiomatic.

## Tests

Set `REQUESTER_TEST_BUCKET` variable to name of test bucket and run `make test`.

## Usage

1. Authorize your machine with `gcloud auth login`.
2. Run `make deploy`.
3. Publish a message to `requester-instruction` containing a JSON object with the keys `url` and `bucket`.
