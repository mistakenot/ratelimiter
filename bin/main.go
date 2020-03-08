package main

import "github.com/mistakenot/ratelimiter/cmd"

// Due to the requirements of Cloud Functions, the root dir cannot
//  contain a main package, so we put it here instead.
func main() {
	cmd.Execute()
}
