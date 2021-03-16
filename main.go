package main

import (
	"os"

	"bot/discord"
)

func main() {
	// NOTE: as this is just some waiting function.
	// Anything beyond this code is unreachable.
	discord.Initalize(os.Getenv("DISCORD_AUTH"))
}
