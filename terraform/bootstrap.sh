#!/usr/bin/env bash

# Idempotent script to bootstrap the initial GCP resources required to use this repository and terraform.

set -e

# You need gcloud installed.
which gcloud > /dev/null || {
    echo "Ensure you have the gcloud tool installed."
    exit 1
}

# And the gsutil tool (should be installed with gcloud).
which gsutil > /dev/null || {
    echo "Ensure you have the gsutil tool installed."
    exit 1
}

# And terraform.
which terraform > /dev/null || {
    echo "Ensure you have the terraform tool installed."
    exit 1
}

# You will need to have a pre-existing project.
test -n "$TF_VAR_project_id" || {
    echo "Set the TF_VAR_project_id variable."
    exit 1
}

gcloud projects describe "$TF_VAR_project_id" > /dev/null || {
    echo "Either the project doesn't exist or you don't have the right permissions to access it."
    exit 1
}

# This script should be run from the root of the directory.
test -f ./makefile || {
    echo "Run this script from the root of the directory."
    exit 1
}

TERRAFORM_BUCKET="$TF_VAR_project_id-ratelimiter-terraform"

# A bucket is required for terraform state.
gsutil ls "gs://$TERRAFORM_BUCKET" > /dev/null 2>&1 || {
    gsutil mb "gs://$TERRAFORM_BUCKET"
}

RATELIMITER_REPO="ratelimiter"

# And a repository to store the source code.
gcloud source repos describe "$RATELIMITER_REPO" > /dev/null 2>&1 || {
    gcloud source repos create "$RATELIMITER_REPO" 
    git config --global credential.'https://source.developers.google.com'.helper gcloud.sh
    git remote add google "https://source.developers.google.com/p/$TF_VAR_project_id/r/$RATELIMITER_REPO"
    echo "Push to the repository when you are ready with 'git push google master'. This will trigger a build of the image."
}

# And a trigger
gcloud beta builds triggers describe "$RATELIMITER_REPO" > /dev/null 2>&1 || {
    gcloud beta builds triggers create cloud-source-repositories \
        --build-config ./cloudbuild.yaml \
        --repo "$RATELIMITER_REPO" \
        --description "$RATELIMITER_REPO" \
        --branch-pattern ".*"
}

# And finally, init terraform
test -d ./.terraform || {
    terraform init -backend-config=bucket="$TERRAFORM_BUCKET" ./terraform
}