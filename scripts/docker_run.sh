#!/bin/sh
# File: scripts/docker_run.sh
# Purpose: Run the bot locally in Docker with .env file.
# Usage: sh scripts/docker_run.sh
set -e

docker build -t tu-tien-bot .
docker run --rm \
  --env-file .env \
  -p 8080:8080 \
  --name tu-tien-bot \
  tu-tien-bot
