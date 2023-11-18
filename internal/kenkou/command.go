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
				Name:        "alarm",
				Description: "Configure alarm ON/OFF. If not specified, returns the current setting.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "active",
						Description: "T/F",
						Required:    false,
					},
				},
			},
			{
				Name:        "weekday",
				Description: "Configure Weekday-Only setting. If not specified, returns the current setting.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "only",
						Description: "T/F",
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
		case "channel":
			guildId := i.GuildID
			setting, _ := GetGuildKenkouSetting(guildId)

			// 指定がない場合は現在のチャンネルを返す
			if len(subcommand.Options) != 0 {
				newAlarmChannel := subcommand.Options[0].ChannelValue(session).ID
				setting.AlarmChannel = &newAlarmChannel
				setting = SaveGuildKenkouSetting(setting)
			}

			var content string
			if setting.AlarmChannel == nil {
				content = "Kenkou channel is not set.\nPlease use `/setting channel <channel>` command to set channel."
			} else {
				content = "Current kenkou channel is <#" + *setting.AlarmChannel + ">"
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
			setting, _ := GetGuildKenkouSetting(guildId)

			if len(subcommand.Options) != 0 {
				newTimeString := subcommand.Options[0].StringValue()
				newTime, _ := time.Parse(time.TimeOnly, newTimeString+":00")
				setting, _ := GetGuildKenkouSetting(guildId)
				setting.AlarmTime = newTime
				setting = SaveGuildKenkouSetting(setting)
			}

			content := "Current kenkou time is " + setting.AlarmTime.Format("15:04")

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})

		case "alarm":
			guildId := i.GuildID
			setting, _ := GetGuildKenkouSetting(guildId)

			if len(subcommand.Options) != 0 {
				newAlarmActive := subcommand.Options[0].BoolValue()
				setting, _ := GetGuildKenkouSetting(guildId)
				setting.AlarmActive = newAlarmActive
				setting = SaveGuildKenkouSetting(setting)
			}

			content := "Current setting is " + fmt.Sprintf("%t", setting.AlarmActive)

			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})

		case "weekday":
			guildId := i.GuildID
			setting, _ := GetGuildKenkouSetting(guildId)

			if len(subcommand.Options) != 0 {
				newAlarmWeekdayOnly := subcommand.Options[0].BoolValue()
				setting, _ := GetGuildKenkouSetting(guildId)
				setting.AlarmWeekdayOnly = newAlarmWeekdayOnly
				setting = SaveGuildKenkouSetting(setting)
			}

			content := "Current setting is " + fmt.Sprintf("%t", setting.AlarmWeekdayOnly)

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
			content = content + "\n------------------------------------------------------------------------------------------"
			content = content + "\n|        Guild        | AlarmActive |     AlarmChannel    | AlarmTime | AlarmWeekdayOnly |"
			content = content + "\n------------------------------------------------------------------------------------------"
			for _, setting := range kenkouSettings {
				var alarmChannel string
				if setting.AlarmChannel != nil {
					alarmChannel = *setting.AlarmChannel
				} else {
					alarmChannel = "undefined"
				}
				content = content +
					"\n| " + fmt.Sprintf("%19s", setting.Guild) +
					" | " + fmt.Sprintf("%11s", fmt.Sprintf("%t", setting.AlarmActive)) +
					" | " + fmt.Sprintf("%19s", alarmChannel) +
					" | " + fmt.Sprintf("%9s", setting.AlarmTime.Format("15:04")) +
					" | " + fmt.Sprintf("%16s", fmt.Sprintf("%t", setting.AlarmWeekdayOnly)) +
					" |"
			}
			content = content + "\n------------------------------------------------------------------------------------------```"

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
