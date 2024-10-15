#!/bin/bash
set -e

echo "===== Starting the Deployment Process ====="

# Step 1: Git operations
echo "Step 1: Checking and pulling the latest changes from the Git repository..."
git fetch origin

# Force pull changes
echo "Forcing pull of latest changes..."
git reset --hard origin/$(git rev-parse --abbrev-ref HEAD)
if git pull --force; then
    echo "Successfully pulled latest changes (force pull)."
else
    echo "Force pull failed. Exiting."
    exit 1
fi

# Step 2: Check current running version
echo "Step 2: Checking current running version..."
if docker compose ps | grep -q "alias-api"; then
    echo "alias-api is currently running."
else
    echo "alias-api is not currently running."
fi

# Step 3: Build and update
echo "Step 3: Building and updating the Docker service..."
if docker compose up -d --build --no-deps alias-api-app; then
    echo "Successfully built and updated alias-api service."
else
    echo "Failed to build or update alias-api-app service. Rolling back..."
    docker compose up -d --no-deps alias-api-app
    echo "Rolled back to previous version."
    exit 1
fi

# Step 4: Cleanup
echo "Step 4: Cleaning up unused Docker resources..."
docker system prune -f

echo "===== Deployment process completed ====="