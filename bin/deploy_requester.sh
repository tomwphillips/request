#!/bin/bash

set -eu

gcloud functions deploy ConsumePubSub --runtime go111 --trigger-topic requester-instruction
