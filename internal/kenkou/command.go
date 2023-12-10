package kenkou

import (
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
				Name:        "dump",
				Description: "Show current setting.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "edit",
				Description: "Edit setting.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "active",
						Description: "Enable Kenkou Alarm.",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "Kenkou Alarm channel. Must be VoiceChannel",
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildVoice,
						},
						Required: false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "time",
						Description: "Kenkou Alarm time. Must be like 01:00",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "weekday",
						Description: "Enable Kenkou Alarm only on weekdays",
						Required:    false,
					},
				},
			},
			{
				Name:        "dump-all",
				Description: "Dump all kenkou settings",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
}

func CommandHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "kenkou":
		guildId := i.GuildID
		channel, _ := session.Channel(i.ChannelID)
		if channel.Type != 2 { // if specified is not ChannelTypeGuildVoice, change setting.Channel
			setting, _ := GetGuildKenkouSetting(guildId)
			if setting.AlarmChannel == nil {
				session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No channel found.\nPlease use `/setting update-kenkou-channel` command to set channel.",
						Flags:   1 << 6,
					},
				})

				return
			}
			channel, _ = session.Channel(*setting.AlarmChannel)
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
		case "edit":
			options := subcommand.Options
			guildId := i.GuildID
			setting, _ := GetGuildKenkouSetting(guildId)

			for _, option := range options {
				switch option.Name {
				case "active":
					setting.AlarmActive = option.BoolValue()
				case "channel":
					setting.AlarmChannel = &option.ChannelValue(session).ID
				case "time":
					// TODO: 入力形式バリデーション
					newTime, _ := time.Parse(time.TimeOnly, option.StringValue()+":00")
					setting.AlarmTime = newTime
				case "weekday":
					setting.AlarmWeekdayOnly = option.BoolValue()
				}
			}

			SaveGuildKenkouSetting(setting)

			content := GetDumpString([]KenkouSetting{setting})

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})

		case "dump":
			guildId := i.GuildID
			setting, _ := GetGuildKenkouSetting(guildId)

			// TODO: DumpではGuildは出したくない気もする
			content := GetDumpString([]KenkouSetting{setting})

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})

		case "dump-all":
			kenkouSettings, _ := GetKenkouSettings()

			content := GetDumpString(kenkouSettings)

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
