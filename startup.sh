#! /bin/bash
PROJECT_ID="lazuli-venturas"
INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
ZONE=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | awk -F/ '{print $NF}')
CODE_DIRECTORY="ninja-combined-crawler"
export PROJECT_DIR=$INSTANCE_NAME
# Update the instance with shielded-learn-integrity-policy
gcloud compute instances update $INSTANCE_NAME --zone $ZONE --shielded-learn-integrity-policy

# Set ulimit
ulimit -n 1000000

# Stop the MongoDB service
sudo service mongod stop

# Change to the code directory
cd /root/$CODE_DIRECTORY

# Discard all local changes
git reset --hard HEAD
git clean -df

# Checkout the dev branch
git checkout dev

# Pull latest changes from the 'dev' branch of the Git repository
git pull

cp /root/$CODE_DIRECTORY/apps/$INSTANCE_NAME/.env .env

# Download the setup_docker.sh script
curl -O https://raw.githubusercontent.com/uzzalhcse/awesome-bash-scripts/main/setup_docker.sh

# Make the script executable
chmod +x setup_docker.sh

# Run the setup_docker.sh script, passing INSTANCE_NAME as an argument
sudo ./setup_docker.sh

# Build and start the Docker containers
echo "Building and starting Docker containers..."
sudo -E docker compose up --build -d


# Wait for the 'crawler' container to exit
echo "Waiting for 'crawler' to exit..."
sudo docker wait crawler
# Remove the setup_docker.sh script
rm setup_docker.sh

echo "Stopping the instance..."
# Stop the instance
ACCESS_TOKEN=$(gcloud auth print-access-token)
URL="https://compute.googleapis.com/compute/v1/projects/$PROJECT_ID/zones/$ZONE/instances/$INSTANCE_NAME/stop"

curl -s -X POST "$URL" -H "Authorization: Bearer $ACCESS_TOKEN"
