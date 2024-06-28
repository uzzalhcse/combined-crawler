#! /bin/bash

# Set environment variables
export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=0
export DB_USERNAME=lazuli
export DB_PASSWORD=x1RWo6cqFtHiaAHce5HB
export DB_HOST=mongo
export DB_PORT=27017
export USER_AGENT="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
export APP_ENV=production
export API_USERNAME=lazuli
export API_PASSWORD=ninja
export GCP_CREDENTIALS_PATH=/app/gcp-file-upload-key.json
export MONGO_INITDB_ROOT_USERNAME=lazuli
export MONGO_INITDB_ROOT_PASSWORD=x1RWo6cqFtHiaAHce5HB

INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
ZONE=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | awk -F/ '{print $NF}')
CODE_DIRECTORY="ninja-combined-crawler"

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

# Pull latest changes from the 'dev' branch of the Git repository
git pull origin dev

# Download the setup_docker.sh script
curl -O https://raw.githubusercontent.com/uzzalhcse/awesome-bash-scripts/main/setup_docker.sh

# Make the script executable
chmod +x setup_docker.sh

# Run the setup_docker.sh script, passing INSTANCE_NAME as an argument
sudo ./setup_docker.sh $INSTANCE_NAME

# Remove the setup_docker.sh script
rm setup_docker.sh
