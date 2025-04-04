#!/bin/bash
#
# Slack Standup Updater Wrapper Script
# This allows running the Go program without building it

# Change to the script's directory
cd "$(dirname "$0")"

# Run the Go program directly
go run main.go

# Exit with the same status as the Go program
exit $? 