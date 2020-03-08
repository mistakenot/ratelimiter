variable "project_id" {
    type = "string"
    description = "GCP project id."
}

variable "region" {
    type = "string"
    description = "GCP region to deploy into."
    default = "us-central1"
}

variable "max_requests_in_period" {
    description = "Max requests the limiter will accept from a single id, for each time window."
    default = 20
}

variable "period_duration_seconds" {
    description = "Length of time in seconds before the limiter resets an id's request count."
    default = 10
}

variable "source_repo_sha" {
    type = "string"
    description = "Commit hash to deploy."
}