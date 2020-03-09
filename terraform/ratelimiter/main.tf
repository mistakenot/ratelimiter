provider "google" {
  project     = "${var.project_id}"
  region      = "${var.region}"
  version     = "~> 3.11"
}

resource "google_redis_instance" "ratelimiter_redis" {
  name           = "ratelimiter-redis"
  region         = "${var.region}"
  memory_size_gb = 1
}

# Required to access Redis
resource "google_vpc_access_connector" "connector" {
  name          = "ratelimiter"
  region        = "${var.region}"
  ip_cidr_range = "10.8.0.0/28"
  network       = "default"
}

resource "google_cloudfunctions_function" "ratelimiter" {
  name        = "ratelimiter"
  description = "Ratelimiter service."
  runtime     = "go113"
  entry_point = "Index"
  trigger_http = true
  source_repository {
    url =   "https://source.developers.google.com/projects/${var.project_id}/repos/ratelimiter/revisions/${var.source_repo_sha}"
  }

  # TODO I had an issue with this field where the deploy completed without errors, but didn't set this option on the function.
  #  I ended up setting it with gcloud.
  vpc_connector = "projects/${var.project_id}/locations/${var.region}/connectors/${google_vpc_access_connector.connector.name}"

  environment_variables = {
    REDIS_URL = "${google_redis_instance.ratelimiter_redis.host}:${google_redis_instance.ratelimiter_redis.port}"
    MAX_REQUESTS_IN_PERIOD = "${var.max_requests_in_period}"
    PERIOD_DURATION_IN_SECONDS = "${var.period_duration_seconds}"
  }
}

# TODO Public access. 
resource "google_cloudfunctions_function_iam_member" "invoker" {
  region         = "${var.region}"
  cloud_function = "${google_cloudfunctions_function.ratelimiter.name}"
  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}
