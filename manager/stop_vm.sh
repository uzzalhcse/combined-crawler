#!/bin/bash

PROJECT_ID="lazuli-venturas-stg"
INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
ZONE=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | awk -F/ '{print $NF}')
ACCESS_TOKEN=$(gcloud auth print-access-token)

# Get the status of the VM
STATUS=$(curl -H "Authorization: Bearer $ACCESS_TOKEN" \
    "https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME" \
    | jq -r .status)

if [ "$STATUS" != "RUNNING" ]; then
    # Start the VM if it is not running
    START_URL="https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME/start"
    curl -X POST $START_URL -H "Authorization: Bearer $ACCESS_TOKEN"
    echo "Instance $INSTANCE_NAME in zone $ZONE has been started."
else
    echo "Instance $INSTANCE_NAME in zone $ZONE is already running."
fi
