package internal

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func PlayHotaru(session *discordgo.Session, guildID string, channelID string) {
	vc, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Fatalf("Can't join vc: %v", err)
	}
	log.Println("|_･) VC Joined.")

	encodeSession, err := dca.EncodeFile("./assets/hotaru.mp3", dca.StdEncodeOptions)
	if err != nil {
		log.Fatalf("Can't encode music: %v", err)
	}

	vc.Speaking(true)

	done := make(chan error)
	dca.NewStream(encodeSession, vc, done)

	err = <-done
	if err != nil && err != io.EOF {
		log.Println("err", err)
	}

	vc.Speaking(false)
}

func ForceKenkou(session *discordgo.Session, guildID string, channelID string) {
	forceFlag := false
	for k, v := range session.VoiceConnections {
		if k == guildID && v.ChannelID == channelID {
			forceFlag = true
		}
	}
	if !forceFlag {
		log.Println(">< I'm not in VC.")
		return
	}

	guild, err := session.Guild(guildID)
	if err != nil {
		log.Fatalf("Can't get guild: %v", err)
		return
	}

	members, _ := session.GuildMembers(guild.ID, "", 1000)
	for _, member := range members {
		go session.GuildMemberMove(guild.ID, member.User.ID, nil)
	}

	log.Println(">< All kicked.")
	return
}

func CheckWeekday(time time.Time) bool {
	url := "https://s-proj.com/utils/checkHoliday.php?kind=h&date=" + time.Format("20060102")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray) == "else"
}
