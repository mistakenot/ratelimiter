terraform {
  backend "gcs" {
      # Config is passed in when the init command is called.
  }
}

module "ratelimiter" {
    source = "./ratelimiter"
    project_id = "${var.project_id}"
    region = "us-central1"
    source_repo_sha = "${var.source_repo_sha}"
}