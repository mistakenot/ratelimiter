package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ratelimiter",
	Short: "A simple rate limiter service.",
	Long:  "A simple, distributed rate limiter service built using Go and Redis.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//  Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
