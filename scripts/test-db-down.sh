#!/bin/bash

echo "Removing db containers"
docker compose -f ./docker/db.docker-compose.yaml down