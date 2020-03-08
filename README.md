# Rate limiter

A simple, distributed rate limiter service built with Go, Redis, Terraform and GCP.

# Usage

The service exposes a REST api described in the [Openapi](./api/openapi.yaml) document. It exposes the following endpoints:

- `GET /` and `GET /healthz` perform a health check and returns `200` if the service can connect to Redis.
- `GET /token/{token_id}` will increment the request count for the given token, and will return a response that details the number of tokens remaining, and the number of seconds until reset.

# Features

- [Infrastructure](./terraform) described with terraform, including a [bootstrap script](./terraform/bootstrap.sh), for deploying to GCP.
- CI/CD config in [cloudbuild.yaml](./cloudbuild.yaml).

# Developers

- `make build` builds the binary.
- `make test` runs unit tests. Note that this requires an instance of Redis 4 to be running locally.
- `make up` starts a full local environment with docker-compose.
- `make deploy` runs the cloud init and deployment scripts, powered by Terraform.

# Todo

- There are a few `TODO`'s in the code base for small code improvements. (`grep -r TODO .`)
- Functions are deploy to public endpoints, without IAM authentication.