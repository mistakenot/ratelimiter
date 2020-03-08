build:
	go build ./ratelimiter.go

test:
	go test ./pkg/limiter

# These tests require a redis instance to be running locally.
test-full: test
	go test ./internal/redis

up:
	docker-commpose down && \
	docker-compose up -d