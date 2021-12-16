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

	case "report":
		fmt.Printf("%#v\n", args)

		// Check for empty reports.
		if len(args) == 1 {
			s.ChannelMessageSend(m.ChannelID, "Please fill out information about your report.")
			break
		}

		// Delete the report message.
		if err = s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			break
		}

		content := strings.Join(args[1:], " ")
		report := fmt.Sprintf("Report by %v:\n%v", m.Author.Mention(), content)
		_, err = s.ChannelMessageSend(utils.ReportChannelID, report)
		if err != nil {
			break
		}

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

		gifURL, errGo := gif.GetRandomGif("nervous")
		if errGo != nil {
			log.Println(errGo)
		}
		_, errGo = s.ChannelMessageSend(m.ChannelID, "I'm nervous REEE, "+gifURL)
		if errGo != nil {
			log.Println(errGo)
		}

		deployCommand := exec.Command("usw", "deploy", args[1])
		err = deployCommand.Start()
		if err != nil {
			log.Println("Couldn't deploy:", err)
			break
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Successfully deployed the branch \""+args[1]+"\"!")
		if err != nil {
			break
		}

		// Send a succesfull deploy gif.
		gifURL, err = gif.GetRandomGif("success")
		if err != nil {
			break
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "I'm happy REEE, "+gifURL)
	}

	if err != nil {
		log.Println(err)
	}
}
