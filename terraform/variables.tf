variable "project_id" {
    type = string
    description = "GCP project id."
}

variable "region" {
    type = string
    description = "GCP region to deploy into."
    default = "us-central1"
}

variable "source_repo_sha" {
    type = string
    description = "Commit hash to deploy to functions."
}

variable "environment" {
    type = string
    description = "The environment that is being deployed to, e.g. development."
    default = "development"
}