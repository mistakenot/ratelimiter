package main

import (
	"fmt"
	"os"

	"datadyne.io/ratelimiter/cmd/ratelimiter"
)

func main() {
	if err := ratelimiter.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
