terraform {
  backend "gcs" {
      # Config is passed in when the init command is called.
  }
}

data "google_client_openid_userinfo" "me" {
}

module "ratelimiter" {
  source = "./ratelimiter"
  project_id = var.project_id
  region = var.region
  source_repo_sha = var.source_repo_sha
  deployer_label = data.google_client_openid_userinfo.me.email
  environment_label = var.environment
}