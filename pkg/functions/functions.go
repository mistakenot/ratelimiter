package functions

import (
	"net/http"

	"github.com/mistakenot/ratelimiter/internal/serve"
	"github.com/mistakenot/ratelimiter/pkg/limiter"
	"github.com/spf13/viper"
)

// Used as entrypoint for google functions
var Index http.HandlerFunc

func init() {
	viper.AutomaticEnv()

	redisUrl := viper.GetString("REDIS_URL")
	maxRequestsInPeriod := viper.GetInt("MAX_REQUESTS_IN_PERIOD")
	periodDurationInSeconds := viper.GetInt("PERIOD_DURATION_IN_SECONDS")
	config := limiter.RateLimiterConfig{maxRequestsInPeriod, periodDurationInSeconds, redisUrl}

	limiter, err := limiter.Create(config)

	if err != nil {
		panic(err)
	}

	Index = serve.CreateHandler(limiter).ServeHTTP
}
