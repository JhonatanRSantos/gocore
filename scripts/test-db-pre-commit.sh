#!/bin/bash

# Temporary fix

echo "Creating db containers"
docker compose -f ./docker/db.docker-compose.yaml up -d

sleep 30