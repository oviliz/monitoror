#!/usr/bin/env bash
# Do not use this script manually, Use makefile

set -e

source ./scripts/setup-variables.sh

###############################################
# This script is used to start monitoror core #
###############################################

# Set environment (default: development)
export MO_ENV=${MO_ENV:-$MB_ENVIRONMENT}

# Build and run
# Avoid firewall issue on Windows
source ./scripts/build/build.sh

function removeTarget() {
  rm $target
}
trap removeTarget TERM INT

$target
