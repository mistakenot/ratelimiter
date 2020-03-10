#!/usr/bin/env bash

# This script will bootstrap the initial GCP resources required to use this repository and terraform.

set -e

# Requires terraform.
which terraform > /dev/null || {
    echo "Ensure you have the terraform tool installed."
    exit 1
}

# This script should be run from the root of the directory.
test -f ./makefile || {
    echo "Run this script from the root of the directory."
    exit 1
}

# Environment variable allows us to use tfvar files to differentiate between different environments.
test -n "$ENVIRONMENT" || {
    echo "Set the ENVIRONMENT variable to correspond to a file in ./terraform/environments."
    exit 1
}

ENVIRONMENT_FILE="terraform/environments/$ENVIRONMENT.tfvars"

test -f "$ENVIRONMENT_FILE" || {
    echo "Cannot find file $ENVIRONMENT_FILE."
    echo "Please ensure that the ENVIRONMENT variable is correct."
    exit 1
}

terraform apply -var-file="$ENVIRONMENT_FILE" ./terraform