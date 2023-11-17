package kenkou

import (
	"fmt"
	"time"

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
}

func SlashCommandHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "kenkou":
		guildId := i.GuildID
		channel, _ := session.Channel(i.ChannelID)
		if channel.Type != 2 { // if specified is not ChannelTypeGuildVoice, change setting.Channel
			setting, _ := GetGuildKenkouSetting(guildId)
			if setting.ChannelId == nil {
				session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No channel found.\nPlease use `/setting update-kenkou-channel` command to set channel.",
						Flags:   1 << 6,
					},
				})

				return
			}
			channel, _ = session.Channel(*setting.ChannelId)
		}
		go ForceKenkou(session, guildId, channel.ID)

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
				setting, _ := GetGuildKenkouSetting(guildId)
				if setting.ChannelId == nil {
					content = "Kenkou channel is not set.\nPlease use `/setting channel <channel>` command to set channel."
				} else {
					content = "Current kenkou channel is <#" + *setting.ChannelId + ">"
				}
			} else {
				newChannel := subcommand.Options[0].ChannelValue(session)
				newChannelId := newChannel.ID
				setting, _ := GetGuildKenkouSetting(guildId)
				setting.ChannelId = &newChannelId
				UpdateGuildKenkouSetting(setting)
				content = "Current kenkou channel is <#" + *setting.ChannelId + ">"
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
				setting, _ := GetGuildKenkouSetting(guildId)
				content = "Current kenkou time is " + setting.Time.Format("15:04")
			} else {
				newTimeString := subcommand.Options[0].StringValue()
				newTime, _ := time.Parse(time.TimeOnly, newTimeString+":00")
				setting, _ := GetGuildKenkouSetting(guildId)
				setting.Time = newTime
				UpdateGuildKenkouSetting(setting)
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
			kenkouSettings, _ := GetKenkouSettings()
			content := "**Kenkou settings**"
			content = content + "\n```"
			content = content + "\n--------------------------------------------------------"
			content = content + "\n|       Guild ID       |      Channel ID      |  TIME  |"
			content = content + "\n--------------------------------------------------------"
			for _, setting := range kenkouSettings {
				var channelId string
				if setting.ChannelId != nil {
					channelId = *setting.ChannelId
				} else {
					channelId = "undefined"
				}
				content = content + "\n| " + fmt.Sprintf("%20s", setting.GuildId) + " | " + fmt.Sprintf("%20s", channelId) + " | " + fmt.Sprintf("%6s", setting.Time.Format("15:04")) + " |"
			}
			content = content + "\n--------------------------------------------------------```"

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})
		}
	}
}
