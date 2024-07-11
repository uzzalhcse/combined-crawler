#!/bin/bash

PROJECT_ID="lazuli-venturas"
INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
ZONE=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | awk -F/ '{print $NF}')
ACCESS_TOKEN=$(gcloud auth print-access-token)

URL="https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME/stop"

curl -X POST $URL -H "Authorization: Bearer $ACCESS_TOKEN"
