package cmd

import (
	"fmt"
	"os"

	"github.com/mistakenot/ratelimiter/internal/serve"
	"github.com/mistakenot/ratelimiter/pkg/limiter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts the service.",
	Run: func(cmd *cobra.Command, args []string) {
		redisUrl, _ := cmd.Flags().GetString("redis-url")

		if redisUrl == "" {
			redisUrl = "localhost:6379"
		}

		port, err := cmd.Flags().GetInt("port")

		if port == 0 {
			port = 8080
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		maxRequestsInPeriod, err := cmd.Flags().GetInt("max-requests-in-period")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if maxRequestsInPeriod == 0 {
			fmt.Print("--max-requests-in-period required.")
			os.Exit(1)
		}

		periodDurationInSeconds, err := cmd.Flags().GetInt("period-duration-seconds")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if periodDurationInSeconds == 0 {
			fmt.Print("--period-duration-seconds required.")
			os.Exit(1)
		}

		config := limiter.RateLimiterConfig{maxRequestsInPeriod, periodDurationInSeconds, redisUrl}
		limiter, err := limiter.Create(config)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := serve.Serve(limiter, port); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCommand)
	startCommand.Flags().IntP("port", "p", viper.GetInt("PORT"), "Port to listen on. [defaults to 8080]")
	startCommand.Flags().StringP("redis-url", "r", viper.GetString("REDIS_URL"), "Redis instance connection url. [defaults to localhost:6379]")
	startCommand.Flags().Int("max-requests-in-period", viper.GetInt("MAX_REQUESTS_IN_PERIOD"), "Max number of requests in each period, per user.")
	startCommand.Flags().Int("period-duration-seconds", viper.GetInt("PEROID_DURATION_SECONDS"), "Length of each measured period, in seconds.")
}
