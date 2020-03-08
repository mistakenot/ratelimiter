#!/usr/bin/env bash

# This script will bootstrap the initial GCP resources required to use this repository and terraform.

set -e

# Required terraform.
which terraform > /dev/null || {
    echo "Ensure you have the terraform tool installed."
    exit 1
}

# This script should be run from the root of the directory.
test -f ./makefile || {
    echo "Run this script from the root of the directory."
    exit 1
}

terraform apply ./terraform