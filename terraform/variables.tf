variable "project_id" {
    type = "string"
    description = "GCP project id."
}

variable "region" {
    type = "string"
    description = "GCP region to deploy into."
    default = "us-central1"
}