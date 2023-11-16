package kenkou

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var HotaruDCABuffer = make([][]byte, 0)

func KenkouBatch(session *discordgo.Session) {
	nowTime := time.Now()
	if !checkWeekday(nowTime) {
		return
	}
	kenkouSettings, _ := GetKenkouSettings()
	for _, setting := range kenkouSettings {
		kenkouTime := setting.Time
		kenkouTime = kenkouTime.Add(time.Minute * -5)
		// MEMO: UTCとJSTごっちゃだけどなんか動くのでヨシッ！
		if nowTime.Hour() == kenkouTime.Hour() && nowTime.Minute() == kenkouTime.Minute() {
			go ForceKenkou(session, setting.GuildId, setting.ChannelId)
		}
	}
}

func ForceKenkou(s *discordgo.Session, guildID string, channelID string) {
	channel, _ := s.Channel(channelID)
	if channel.Type != 2 { // is not ChannelTypeGuildVoice
		channelID = os.Getenv("DEFAULT_CHANNEL_ID")
	}

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
