#!/usr/bin/env bash

CONTAINER_RUNTIME=$(command -v podman &> /dev/null && echo podman || echo docker)

mkdir -p "/tmp/datacoviddb"

$CONTAINER_RUNTIME run \
--rm \
--name japanvisacovidbotdb \
-e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
-e POSTGRES_USER="${POSTGRES_USER}" \
-e POSTGRES_DB="${POSTGRES_DB}" \
-p 5432:5432 \
-v ./sql/tables.sql:/docker-entrypoint-initdb.d/1.sql \
-v /tmp/datacoviddb:/var/lib/postgresql/data \
docker.io/library/postgres:14.1