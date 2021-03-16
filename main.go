package main

import (
	"os"

	"bot/discord"
)

func main() {
	discord.Initalize(os.Getenv("DISCORD_AUTH"))
}
