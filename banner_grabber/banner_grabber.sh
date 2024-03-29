#!/bin/bash

# Assign the first script argument to a variable
IP=$1

# Check if an IP address is provided
if [ -z "$IP" ]; then
  echo "Usage: $0 <IP>"
  exit 1
fi

# Run wget command with the provided IP and return the output
output=$(wget "$IP" -q -S 2>&1)

echo "Command Output:"
echo "$output"
