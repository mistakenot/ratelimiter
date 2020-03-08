export TF_VAR_project_id ?= dd-dev-exam

build:
	go build ./bin/main.go

# These tests require a redis instance to be running locally.
test:
	go test ./pkg/limiter

down: 
	docker-compose down

up: down
	docker-compose up -d

deploy:
	bash ./terraform/bootstrap.sh && bash ./terraform/deploy.sh
