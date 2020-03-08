#!/usr/bin/env bash

# This script will bootstrap the initial GCP resources required to use this repository and terraform.

set -e

# Required terraform.
which terraform > /dev/null || {
    echo "Ensure you have the terraform tool installed."
    exit 1
}