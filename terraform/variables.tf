variable "project_id" {
    type = "string"
    description = "GCP project id."
}

variable "region" {
    type = "string"
    description = "GCP region to deploy into."
    default = "us-central1"
}

variable "image_version" {
    type = "string"
    description = "The version of the service image. Defaults to latest."
    default = "latest"
}