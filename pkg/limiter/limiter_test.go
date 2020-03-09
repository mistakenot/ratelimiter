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
	bucketStart, bucketEnd := calculateCurrentBucketStartAndEnd(now, 10)
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

	// TODO this isn't ideal as it relies on the clock and could be unstable.
	// Possible way to resolve would be to inject a "GetTime" service of some sort
	//  so that you can predictably set it for testing. Out of scope at the mo.
	for i := 1; i < 6; i++ {
		result, err := limiter.TakeToken(key)

		if err != nil {
			t.Error(err)
		}

		if result.RequestsRemaining != (5 - i) {
			t.Errorf("Expected %v, got %v", (5 - i), result.RequestsRemaining)
		}

		if result.SecondsToReset > 100 || result.SecondsToReset < 0 {
			// TODO this could be better asserted.
			t.Errorf("SecondsToReset is invalid: %v", result.SecondsToReset)
		}
	}

	// Should be zero there after
	result, err := limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}

	if result.RequestsRemaining != 0 {
		t.Errorf("Expected %v, got %v", 0, result.RequestsRemaining)
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
	if result.RequestsRemaining != 1 {
		t.Errorf("Expected %v, got %v", 1, result.RequestsRemaining)
	}

	// Take the second token
	result, err = limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}
	if result.RequestsRemaining != 0 {
		t.Errorf("Expected %v, got %v", 0, result.RequestsRemaining)
	}

	// Wait for reset.
	// TODO again, check if this could cause test instability.
	time.Sleep(1 * time.Second)

	result, err = limiter.TakeToken(key)

	if err != nil {
		t.Error(err)
	}

	if result.RequestsRemaining != 1 {
		t.Errorf("Expected %v, got %v", 1, result.RequestsRemaining)
	}
}
