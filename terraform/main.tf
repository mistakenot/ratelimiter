terraform {
  backend "gcs" {
    bucket  = "${var.project_id}-ratelimiter-terraform"
  }
}

module "ratelimiter" {
    source = "./ratelimiter"
    project = "${var.project_id}"
    region = "us-central1"
}