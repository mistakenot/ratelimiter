package limiter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestCreateFailsWithInvalidArgs(t *testing.T) {
	_, err := Create(RateLimiterConfig{MaxRequestsInPeriod: -1, PeriodDurationSeconds: 1000, RedisConnectionString: "localhost:1234"})
	if err == nil {
		t.Error("Didn't return an error.")
	}
	_, err = Create(RateLimiterConfig{MaxRequestsInPeriod: 1000, PeriodDurationSeconds: 0, RedisConnectionString: "localhost:1234"})
	if err == nil {
		t.Error("Didn't return an error.")
	}
	_, err = Create(RateLimiterConfig{MaxRequestsInPeriod: 1000, PeriodDurationSeconds: -1, RedisConnectionString: "localhost:1234"})
	if err == nil {
		t.Error("Didn't return an error.")
	}
}

func TestCalculateBucketStats(t *testing.T) {
	now := time.Now().Unix()
	bucketStart, bucketEnd := CalculateCurrentBucketStartAndEnd(now, 10)
	expectedBucketStart := now - (now % 10)
	expectedBucketEnd := expectedBucketStart + 10

	t.Logf("Test epoch is %v", now)

	if bucketStart != expectedBucketStart {
		t.Errorf("Expected %v, got %v", expectedBucketStart, bucketStart)
	}

	if bucketEnd != expectedBucketEnd {
		t.Errorf("Expected %v, got %v", expectedBucketEnd, bucketEnd)
	}
}

func TestTakeTokenFreshKey(t *testing.T) {
	key := fmt.Sprint(rand.Int())
	config := RateLimiterConfig{5, 100, "localhost:6379"}
	limiter, err := Create(config)

	if err != nil {
		t.Error(err)
	}

	// TODO this isn't ideal as it relies on the clock and will be unstable.
	for i := 1; i < 6; i++ {
		result, err := limiter.TakeToken(key)

		if err != nil {
			t.Error(err)
		}

		if result.TicketsRemaining != (5 - i) {
			t.Errorf("Expected %v, got %v", (5 - i), result.TicketsRemaining)
		}

		if result.SecondsToReset > 1000 || result.SecondsToReset < 0 {
			// TODO this could be better asserted.
			t.Errorf("SecondsToReset is invalid: %v", result.SecondsToReset)
		}
	}

	// Should be zero there after
	result, err := limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}

	if result.TicketsRemaining != 0 {
		t.Errorf("Expected %v, got %v", 0, result.TicketsRemaining)
	}
}

func TestTakeTokenFreshKeyAfterExhaustion(t *testing.T) {
	key := fmt.Sprint(rand.Int())
	config := RateLimiterConfig{2, 1, "localhost:6379"}
	limiter, err := Create(config)

	if err != nil {
		t.Error(err)
	}

	// Take the first token
	result, err := limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}
	if result.TicketsRemaining != 1 {
		t.Errorf("Expected %v, got %v", 1, result.TicketsRemaining)
	}

	// Take the second token
	result, err = limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}
	if result.TicketsRemaining != 0 {
		t.Errorf("Expected %v, got %v", 0, result.TicketsRemaining)
	}

	// Wait for reset.
	// TODO again, check if this could cause test instability.
	time.Sleep(1 * time.Second)

	result, err = limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}

	if result.TicketsRemaining != 1 {
		t.Errorf("Expected %v, got %v", 1, result.TicketsRemaining)
	}
}
