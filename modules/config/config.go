package config

import (
	"fmt"
	"os"
)

var (
	ServerPort  = getEnv("SERVER_PORT", "3000")
	DiscordAuth = getEnv("DISCORD_AUTH", "none")
)

func getEnv(name, fallback string) string {
	if val, set := os.LookupEnv(name); set {
		return val
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf(`Env variable not found: %v`, name))
}
