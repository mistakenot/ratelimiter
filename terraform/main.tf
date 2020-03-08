terraform {
  backend "gcs" {
      # Config is passed in when the init command is called.
  }
}

module "ratelimiter" {
    source = "./ratelimiter"
    project_id = "${var.project_id}"
    region = "us-central1"
    image = "gcr.io/${var.project_id}/ratelimiter:${var.image_version}"
}