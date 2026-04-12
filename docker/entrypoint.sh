#!/bin/bash
# Entrypoint script for Coddy sandbox containers

set -e

# Ensure workspace directory exists
mkdir -p /home/coddy/workspace

# Set up environment
export PYTHONDONTWRITEBYTECODE=1
export PYTHONUNBUFFERED=1

# Change to workspace
cd /home/coddy/workspace

# Execute the provided command, or keep container alive
if [ $# -eq 0 ]; then
    # No command provided, keep container running
    exec tail -f /dev/null
else
    # Execute provided command
    exec "$@"
fi
