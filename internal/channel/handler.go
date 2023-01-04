package channel

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ChannelUpdateHandler(session *discordgo.Session, m *discordgo.ChannelUpdate) {
	channel := m.Channel
	var role *discordgo.Role

	guildRoles, _ := session.GuildRoles(m.GuildID)
	for _, v := range guildRoles {
		if v.Name == channel.Name {
			role = v
		}
	}

	if role == nil {
		role, _ = session.GuildRoleCreate(m.GuildID)
		session.GuildRoleEdit(m.GuildID, role.ID, channel.Name, 0, false, 0, true)
	}

	session.ChannelMessageSend(channel.ID, ":bulb: チャンネル名が「"+role.Mention()+"」に変わったよ")
	log.Println("^^ Channel name has changed to " + role.Name)
}
