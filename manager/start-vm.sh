#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 INSTANCE_NAME"
  exit 1
fi

PROJECT_ID="lazuli-venturas"
INSTANCE_NAME=$1 # Take instance name as an argument
ACCESS_TOKEN=$(gcloud auth print-access-token)

# Check if jq is installed
if ! command -v jq &> /dev/null; then
  echo "jq command not found. Installing jq to proceed."
  sudo apt-get install -y jq
fi

# Retrieve the zone of the instance
ZONE=$(gcloud compute instances list --filter="name=($INSTANCE_NAME)" --format="value(zone)")
echo "Instance located at $ZONE Zone."
if [ -z "$ZONE" ]; then
  echo "Instance $INSTANCE_NAME not found."
  exit 1
fi

# Get the status of the VM
STATUS=$(curl -s -H "Authorization: Bearer $ACCESS_TOKEN" \
    "https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME" \
    | jq -r .status)

if [ "$STATUS" != "RUNNING" ]; then
    # Start the VM if it is not running
    START_URL="https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME/start"
    curl -s -X POST $START_URL -H "Authorization: Bearer $ACCESS_TOKEN"
    echo "Instance $INSTANCE_NAME in zone $ZONE has been started."
else
    echo "Instance $INSTANCE_NAME in zone $ZONE is already running."
fi
