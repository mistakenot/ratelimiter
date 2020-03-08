package functions

import (
	"net/http"

	"github.com/mistakenot/ratelimiter/internal/serve"
	"github.com/mistakenot/ratelimiter/pkg/limiter"
	"github.com/spf13/viper"
)

// Used as entrypoint for google functions. It needs to be in the root of the dir.
// Executable is in ./bin
func Index(w http.ResponseWriter, r *http.Request) {
	viper.AutomaticEnv()

	redisUrl := viper.GetString("REDIS_URL")
	maxRequestsInPeriod := viper.GetInt("MAX_REQUESTS_IN_PERIOD")
	periodDurationInSeconds := viper.GetInt("PERIOD_DURATION_IN_SECONDS")
	config := limiter.RateLimiterConfig{maxRequestsInPeriod, periodDurationInSeconds, redisUrl}

	limiter, err := limiter.Create(config)

	if err != nil {
		panic(err)
	}

	serve.CreateHandler(limiter).ServeHTTP(w, r)
}
