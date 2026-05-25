#!/bin/sh
# File: scripts/render_start.sh
# Purpose: Start script for Render deployment.
# Render sets environment variables directly — no .env file needed in production.
set -e

echo "Starting Tu Tien Bot on Render..."
exec ./bot
