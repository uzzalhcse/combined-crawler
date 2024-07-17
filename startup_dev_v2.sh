#! /bin/bash
PROJECT_ID="lazuli-venturas"
INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
ZONE=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | awk -F/ '{print $NF}')
SiteID="sumitool"
# Update the instance with shielded-learn-integrity-policy
gcloud compute instances update $INSTANCE_NAME --zone $ZONE --shielded-learn-integrity-policy

# Set ulimit
ulimit -n 1000000

# Change to the code directory
cd /root
curl -O http://35.243.109.168:8080/binary/$SiteID
chmod +x $SiteID
sudo ./$SiteID

# Download the setup_docker.sh script
curl -s http://35.243.109.168:8080/api/site-secret/$SiteID  > .env

echo "Stopping the instance..."
# Stop the instance
ACCESS_TOKEN=$(gcloud auth print-access-token)
URL="https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME/stop"

curl -s -X POST "$URL" -H "Authorization: Bearer $ACCESS_TOKEN"
