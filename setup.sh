#!/bin/bash

# Stop all running Docker containers
sudo docker stop $(sudo docker ps -a -q)

# Remove all stopped Docker containers
sudo docker rm $(sudo docker ps -a -q)

# Build Docker images using a specific script
sudo deploy/docker/build-all.sh

# Bring up Docker containers using docker-compose
sudo docker-compose -f deploy/docker/docker-compose.yml --compatibility up -d
