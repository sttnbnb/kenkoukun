package kenkou

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var HotaruDCABuffer = make([][]byte, 0)
var vcStatusMap = map[string]*KenkouUserVoiceChatStatus{}

func KenkouBatch(session *discordgo.Session) {
	nowTime := time.Now()
	weekday := checkWeekday(nowTime)

	kenkouSettings, _ := GetKenkouSettings()
	for _, setting := range kenkouSettings {
		if !setting.AlarmActive {
			continue
		}

		if setting.AlarmChannel == nil {
			continue
		}

		if !weekday && setting.AlarmWeekdayOnly {
			continue
		}

		kenkouTime := setting.AlarmTime
		kenkouTime = kenkouTime.Add(time.Minute * -5)
		// MEMO: UTCとJSTごっちゃだけどなんか動くのでヨシッ！
		if nowTime.Hour() == kenkouTime.Hour() && nowTime.Minute() == kenkouTime.Minute() {
			go ForceKenkou(session, setting.Guild, *setting.AlarmChannel)
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
	log.Println("|_･) Start playing Hotaru.")
	log.Println("     Guild: " + guildID + ", Channel: " + channelID)

	vc.Speaking(true)
	for _, buff := range HotaruDCABuffer {
		vc.OpusSend <- buff
	}
	vc.Speaking(false)

	log.Println("     Finish playing.")

	for userID, status := range vcStatusMap {
		if status.GuildID == guildID && status.ChannelID == channelID {
			s.GuildMemberMove(guildID, userID, nil)
			log.Println("     Member disconnected: " + userID)
		}
	}

	vc.Disconnect()

	log.Println("     All done.")
}

func VoiceStateUpdate(s *discordgo.Session, voiceState *discordgo.VoiceStateUpdate) {
	if voiceState.UserID == s.State.User.ID {
		return
	}

	vcStatusMap[voiceState.UserID] = &KenkouUserVoiceChatStatus{
		GuildID:   voiceState.GuildID,
		ChannelID: voiceState.ChannelID,
	}
}
