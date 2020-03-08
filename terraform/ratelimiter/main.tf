provider "google" {
  project     = "${var.project_id}"
  region      = "${var.region}"
}

resource "google_redis_instance" "ratelimiter_redis" {
  name           = "ratelimiter-redis"
  memory_size_gb = 1
}

resource "google_cloud_run_service" "ratelimiter_run" {
  name     = "ratelimiter-run"
  location = "${var.region}"

  template {
    spec {
      containers {
        image = "${var.image}"
        args = [
            "ratelimiter", "start", 
            "--max-requests-in-period", "${var.max_requests_in_period}", 
            "--period-duration-seconds", "${var.period_duration_seconds}",
            "--redis-url", "${google_redis_instance.ratelimiter_redis.host}:${google_redis_instance.ratelimiter_redis.port}"]
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}
