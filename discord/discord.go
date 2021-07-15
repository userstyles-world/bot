package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"bot/handlers"
	"bot/modules/config"
	"bot/utils"
)

func Initalize() {
	discord, err := discordgo.New("Bot " + config.DiscordAuth)
	if err != nil {
		fmt.Println("Wanted to create a new session, but caught error:", err)
		return
	}

	discord.AddHandler(handlers.OnMessage)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("Wanted to open an connection to discord, but caught error:", err)
		return
	}
	utils.Initalize(discord)

	fmt.Println("USw's Guardian is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
