FROM golang:1.13-alpine

WORKDIR /app

ADD . .

RUN go build ./bin/main.go

RUN mv ./main /usr/bin/ratelimiter && chmod +x /usr/bin/ratelimiter

CMD ["ratelimiter", "--help"]