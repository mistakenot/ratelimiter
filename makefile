export TF_VAR_project_id ?= dd-dev-exam

build:
	go build ./main.go

test:
	go test ./pkg/limiter

# These tests require a redis instance to be running locally.
test-full: test
	go test ./internal/redis

up:
	docker-commpose down && \
	docker-compose up -d

deploy:
	bash ./terraform/bootstrap.sh && bash ./terraform/deploy.sh