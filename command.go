package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "kenkou",
		Description: "Force Kenkou",
	},
	{
		Name:        "rename",
		Description: "Rename channel name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "role",
				Required:    true,
			},
		},
	},
	{
		Name:        "newname",
		Description: "Create new channel name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "New channel name",
				Required:    true,
			},
		},
	},
}

func slashCommandsHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "kenkou":
		go playHotaru(session, i.GuildID, i.ChannelID)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Let's Kenkou!",
				Flags:   1 << 6,
			},
		})

	loop:
		for {
			select {
			case <-time.After(5 * time.Minute):
				forceKenkou(session, i.GuildID, i.ChannelID)
				break loop
			}
		}

	case "rename":
		role := i.ApplicationCommandData().Options[0].RoleValue(session, i.GuildID)
		session.ChannelEdit(i.ChannelID, role.Name)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Channel name has changed to" + role.Name,
				Flags:   1 << 6,
			},
		})

	case "newname":
		name := i.ApplicationCommandData().Options[0].StringValue()
		session.ChannelEdit(i.ChannelID, name)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Channel name has changed to" + name,
				Flags:   1 << 6,
			},
		})
	}
}
