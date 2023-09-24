#!/bin/bash

echo "Creating db containers"
docker compose -f ./docker/db.docker-compose.yaml up -d

echo
CONNECTION_ATTEMPT=1
until docker exec pg-test-db psql -U admin -c "select 1" -d test-db > /dev/null 2>/dev/null; do
    echo "wating pg-test-db start. Attempt: $CONNECTION_ATTEMPT"
    CONNECTION_ATTEMPT=$(expr $CONNECTION_ATTEMPT + 1)
    sleep 1
done

echo
CONNECTION_ATTEMPT=1
until docker exec -it mysql-test-db mysql -u admin -pqwerty -e "SELECT 1" > /dev/null 2>/dev/null; do
    echo "wating mysql-test-db start. Attempt: $CONNECTION_ATTEMPT"
    CONNECTION_ATTEMPT=$(expr $CONNECTION_ATTEMPT + 1)
    sleep 1
done

echo
echo "pg-test-db OK"
echo "pg-test-db OK"
echo