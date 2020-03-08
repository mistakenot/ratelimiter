variable "project_id" {
    type = "string"
    description = "GCP project id."
}

variable "region" {
    type = "string"
    description = "GCP region to deploy into."
    default = "us-central1"
}

variable "image" {
    type = "string"
    description = "Ratelimiter full docker tag."
    default = "gcr.io/${var.project_id}/ratelimiter:latest"
}

variable "max_requests_in_period" {
    type = "int"
    description = "Max requests the limiter will accept from a single id, for each time window."
    default = 10
}

variable "period_duration_seconds" {
    type = "int"
    description = "Length of time in seconds before the limiter resets an id's request count."
    default = 1
}