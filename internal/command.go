package internal

import (
	"github.com/shmn7iii/kenkoukun/internal/kenkou"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
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
				Description: "Role name",
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

func SlashCommandsHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "kenkou":
		go kenkou.ForceKenkou(session, i.GuildID, i.ChannelID)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Let's Kenkou!",
				Flags:   1 << 6,
			},
		})

	case "rename":
		role := i.ApplicationCommandData().Options[0].RoleValue(session, i.GuildID)
		session.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
			Name: role.Name,
		})

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Channel name has changed to " + role.Name,
				Flags:   1 << 6,
			},
		})

	case "newname":
		name := i.ApplicationCommandData().Options[0].StringValue()
		session.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
			Name: name,
		})

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "New role has created and channel name has changed to " + name,
				Flags:   1 << 6,
			},
		})
	}
}
