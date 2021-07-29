package handlers

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"bot/modules/gif"
	"bot/utils"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, utils.Prefix) {
		return
	}
	content := strings.TrimPrefix(m.Content, utils.Prefix)
	args := strings.Split(content, " ")

	var err error
	switch args[0] {
	case "ping":
		_, err = s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "help":
		_, err = s.ChannelMessageSend(m.ChannelID, "`ping` - pings the bot\n`uptime` - shows bot uptime\n`deploy <branch>` deploy the server to that branch's current commit\n`help` - shows this message")
	case "uptime":
		uptime := time.Since(utils.LastUptime).Round(time.Second).String()
		if utils.IsDown {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server has been **offline** for %s.", uptime))
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server has been **online** for %s.", uptime))
		}
	case "deploy":
		if len(args) != 2 {
			_, err = s.ChannelMessageSend(m.ChannelID, "Use the syntax: `deploy <branch>`")
			break
		}
		if !utils.ContainsInSlice(m.Member.Roles, utils.AdminRoleID) {
			_, err = s.ChannelMessageSend(m.ChannelID, "You don't have permission to do that.")
			break
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Deploying...")
		if err != nil {
			break
		}

		go func(sesssion *discordgo.Session, channeldID string) {
			// Don't really mind the errors here not being "breaked", it's just a extra.
			gifURL, errGo := gif.GetRandomGif("nervous")
			if errGo != nil {
				log.Println(errGo)
			}
			_, errGo = sesssion.ChannelMessageSend(channeldID, "I'm nervous REEE, "+gifURL)
			if errGo != nil {
				log.Println(errGo)
			}
		}(s, m.ChannelID)

		deployCommand := exec.Command("usw", "deploy", args[1])
		var output []byte
		output, err = deployCommand.Output()
		if err != nil {
			log.Println("Couldn't deploy:", err)
			if err.Error() == "exit status 1" {
				_, err = s.ChannelMessageSend(m.ChannelID, "Couldn't deploy, due to: \""+string(output)+"\"!")
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, "Couldn't deploy, due to: \""+err.Error()+"\"!")
			}
			break
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Successfully deployed the branch \""+args[1]+"\"!")
	}
	if err != nil {
		log.Println(err)
	}
}
