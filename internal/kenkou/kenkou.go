package kenkou

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var HotaruDCABuffer = make([][]byte, 0)

func KenkouBatch(session *discordgo.Session, guildID string, channelID string) {
	nowTime := time.Now()
	if nowTime.Hour() == 0 && nowTime.Minute() == 55 && checkWeekday(nowTime) {
		go ForceKenkou(session, guildID, channelID)
	}
}

func ForceKenkou(s *discordgo.Session, guildID string, channelID string) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Fatalf("Can't join vc: %v", err)
		return
	}
	log.Println("|_･) VC Joined.")

	vc.Speaking(true)
	for _, buff := range HotaruDCABuffer {
		vc.OpusSend <- buff
	}
	vc.Speaking(false)

	members, _ := s.GuildMembers(guildID, "", 1000)
	for _, member := range members {
		go s.GuildMemberMove(guildID, member.User.ID, nil)
	}

	log.Println("(･_| All kicked.")
}
