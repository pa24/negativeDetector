#!/bin/bash

# Pull the latest code from GitHub

git pull origin main

# Build the Docker image
docker build -t negative-detector .

# Stop and remove the old container (if exists)
docker stop negative-detector || true
docker rm negative-detector || true

if [ $(docker ps -aq -f name=negative-detector) ]; then
  docker rm -f negative-detector
fi

# Run the new container
sudo docker run -d -p 8080:8080 --name negative-detector --env-file ../.env negative-detector