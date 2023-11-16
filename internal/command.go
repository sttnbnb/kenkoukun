package internal

import (
	"time"

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
				Name:        "time",
				Description: "Configure kenkou time. If not specified, returns the current time.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "time",
						Description: "Must be like 01:00",
						Required:    false,
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
				setting, _ := kenkou.GetGuildKenkouSetting(guildId)
				if setting.ChannelId == "" {
					content = "Kenkou channel is not set.\nPlease use `/setting channel <channel>` command to set channel."
				} else {
					content = "Current kenkou channel is <#" + setting.ChannelId + ">"
				}
			} else {
				channel := subcommand.Options[0].ChannelValue(session)
				oldSetting, _ := kenkou.GetGuildKenkouSetting(guildId)
				setting := kenkou.NewKenkouSetting(guildId, channel.ID, oldSetting.Time)
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

		case "time":
			guildId := i.GuildID
			content := ""
			if len(subcommand.Options) == 0 {
				// 指定がない場合は現在の設定時刻を返す
				setting, _ := kenkou.GetGuildKenkouSetting(guildId)
				if setting.Time.Unix() == -62135596800 { // 中身ないとこれ
					content = "Kenkou time is not set.\nPlease use `/setting time <01:00>` command to set time."
				} else {
					content = "Current kenkou time is " + setting.Time.Format("15:04")
				}
			} else {
				timeString := subcommand.Options[0].StringValue()
				oldSetting, _ := kenkou.GetGuildKenkouSetting(guildId)
				time, _ := time.Parse(time.TimeOnly, timeString+":00")
				setting := kenkou.NewKenkouSetting(guildId, oldSetting.ChannelId, time)
				kenkou.UpdateGuildKenkouSetting(setting)
				content = "Current kenkou time is " + setting.Time.Format("15:04")
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
			content := "**Kenkou settings**"
			content = content + "\n```     Guild ID     :  Kenkou Channel ID / TIME "
			for _, setting := range KenkouSettings {
				content = content + "\n" + setting.GuildId + ": " + setting.ChannelId + " / " + setting.Time.Format("15:04")
			}
			content = content + "```"

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
