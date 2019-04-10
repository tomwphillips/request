#!/bin/bash

gcloud pubsub topics publish requester-instruction --message "{\"url\": \"http://google.com\"}"
