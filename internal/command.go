package internal

import (
	"fmt"
	"strings"

	"github.com/shmn7iii/kenkoukun/internal/kenkou"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "kenkou",
		Description: "Force Kenkou",
	},
	{
		Name:        "setting",
		Description: "Setting command for Kenkoukun",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channel",
				Description: "Configure kenkou channel. If not specified, returns the current channel.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "Must be VoiceChannel",
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildVoice,
						},
						Required: false,
					},
				},
			},
			{
				Name:        "dump-all-kenkou-settings",
				Description: "Dump all kenkou settings",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
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

func SlashCommandHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "kenkou":
		guildId := i.GuildID
		channel, _ := session.Channel(i.ChannelID)
		if channel.Type != 2 { // is not ChannelTypeGuildVoice
			setting, err := kenkou.GetGuildKenkouSetting(guildId)
			if err != nil {
				session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No channel found.\nPlease use `/setting update-kenkou-channel` command to set channel.",
						Flags:   1 << 6,
					},
				})

				return
			}
			channel, _ = session.Channel(setting.ChannelId)
		}
		go kenkou.ForceKenkou(session, guildId, channel.ID)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Let's Kenkou!",
				Flags:   1 << 6,
			},
		})

	case "setting":
		subcommand := i.ApplicationCommandData().Options[0]
		switch subcommand.Name {
		case "channel":
			guildId := i.GuildID
			content := ""

			if len(subcommand.Options) == 0 {
				// 指定がない場合は現在のチャンネルを返す
				setting, err := kenkou.GetGuildKenkouSetting(guildId)
				content = "Current kenkou channel is <#" + setting.ChannelId + ">"
				if err != nil {
					content = "Kenkou channel is not set.\nPlease use `/setting channel <channel>` command to set channel."
				}
			} else {
				channel := subcommand.Options[0].ChannelValue(session)
				oldSetting, _ := kenkou.GetGuildKenkouSetting(guildId)
				setting := kenkou.KenkouSetting{
					GuildId:   guildId,
					ChannelId: channel.ID,
					Time:      oldSetting.Time,
				}
				kenkou.UpdateGuildKenkouSetting(setting)
				content = "Current kenkou channel is <#" + channel.ID + ">"
			}

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})

		case "dump-all-kenkou-settings":
			KenkouSettings, _ := kenkou.GetKenkouSettings()
			content := "Kenkou settings:"
			for _, setting := range KenkouSettings {
				content = content + "\n  " + setting.GuildId + ": " + setting.ChannelId
			}

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})
		}

	case "rename":
		role := i.ApplicationCommandData().Options[0].RoleValue(session, i.GuildID)
		session.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
			Name: role.Name,
		})

		session.ChannelMessageSend(i.ChannelID, ":bulb: チャンネル名が「"+role.Mention()+"」に変わったよ")

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

		role, _ := session.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
			Name: name,
		})

		session.ChannelMessageSend(i.ChannelID, ":bulb: チャンネル名が「"+role.Mention()+"」に変わったよ")

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "New role has created and channel name has changed to " + name,
				Flags:   1 << 6,
			},
		})
	}
}
