package limiter

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type RateLimiter interface {
	TakeToken(userId string) (RateLimiterResult, error)
	Healthcheck() error
	Close() error
}

type RateLimiterConfig struct {
	MaxRequestsInPeriod   int
	PeriodDurationSeconds int
	RedisConnectionString string
}

type RateLimiterResult struct {
	RequestsRemaining int `json:"requestsRemaining"`
	SecondsToReset    int `json:"secondsToReset"`
}

func Create(config RateLimiterConfig) (RateLimiter, error) {

	if config.MaxRequestsInPeriod < 0 {
		return nil, fmt.Errorf("Config value MaxRequestsInPeriod must be positive integer, got %v", config.MaxRequestsInPeriod)
	}

	if config.PeriodDurationSeconds < 1 {
		return nil, fmt.Errorf("Config value PeriodDurationSeconds must be greater than 1, got %v", config.MaxRequestsInPeriod)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisConnectionString,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return &redisLimiter{client, config}, nil
}

func CalculateCurrentBucketStartAndEnd(now int64, bucketSizeSeconds int64) (int64, int64) {
	bucketStart := now - (now % bucketSizeSeconds)
	bucketEnd := bucketStart + bucketSizeSeconds

	return bucketStart, bucketEnd
}

type redisLimiter struct {
	client *redis.Client
	config RateLimiterConfig
}

func (r *redisLimiter) TakeToken(userId string) (RateLimiterResult, error) {
	// TODO what issues could this cause?
	now := time.Now().Unix()
	bucketStart, bucketEnd := CalculateCurrentBucketStartAndEnd(now, int64(r.config.PeriodDurationSeconds))
	key := fmt.Sprintf("%v:%v", userId, bucketStart)
	expires := time.Duration(bucketEnd-now) * time.Second

	pipe := r.client.TxPipeline()
	incr := pipe.Incr(key)
	pipe.Expire(key, expires)
	_, err := pipe.Exec()

	if err != nil {
		return RateLimiterResult{}, err
	}

	// TODO overflow risk
	ticketsRemaining := r.config.MaxRequestsInPeriod - int(incr.Val())
	secondsToReset := int(bucketEnd - now)

	if ticketsRemaining < 0 {
		ticketsRemaining = 0
	}

	return RateLimiterResult{ticketsRemaining, secondsToReset}, nil
}

func (r *redisLimiter) Close() error {
	return r.client.Close()
}

func (r *redisLimiter) Healthcheck() error {
	_, err := r.client.Ping().Result()
	return err
}
