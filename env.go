package main

import (
	"flag"
	"fmt"
	"os"
)

func parseEnv() {
	if !isFlagSet("redis") {
		*addr = parseStringFromEnv(
			fmt.Sprintf("%s_REDIS", app_env_prefix),
			"localhost:6379")
	}

	if !isFlagSet("inbound") {
		*inbound = parseStringFromEnv(
			fmt.Sprintf("%s_INBOUND", app_env_prefix),
			"irc")
	}

	if !isFlagSet("outbound") {
		*outbound = parseStringFromEnv(
			fmt.Sprintf("%s_OUTBOUND", app_env_prefix),
			app_name)
	}
}

// Parses a string from an env variable and returns it.
func parseStringFromEnv(varName, defaultVal string) string {
	val, set := os.LookupEnv(varName)
	if set {
		return val
	}
	return defaultVal
}

// This is used to check if a flag was set
// Must be called after flag.Parse()
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
