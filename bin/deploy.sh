#!/bin/bash

gcloud functions deploy ConsumePubSub --runtime go111 --trigger-topic requester-instruction

