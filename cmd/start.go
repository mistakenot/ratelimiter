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
			redisUrl = "localhost:5432"
		}

		port, err := cmd.Flags().GetInt("port")

		if port == 0 {
			port = 8080
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		config := limiter.RateLimiterConfig{5, 10, ""}
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
	startCommand.Flags().StringP("redis-url", "r", viper.GetString("NAME"), "Redis instance connection url. [defaults to localhost:5432]")
}
