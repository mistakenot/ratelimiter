export TF_VAR_project_id ?= dd-dev-exam

build:
	go build ./main.go

# These tests require a redis instance to be running locally.
test:
	go test ./internal/redis && \
	go test ./pkg/limiter

up:
	docker-compose down && \
	docker-compose up -d

deploy:
	bash ./terraform/bootstrap.sh && bash ./terraform/deploy.sh