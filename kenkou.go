package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func kenkouBatch(session *discordgo.Session) bool {
	nowTime := time.Now()
	if nowTime.Hour() == 0 && nowTime.Minute() == 55 && checkWeekday(nowTime) {
		go playHotaru(session, DefaultGuildID, DefaultChannelID)
		return true
	} else if nowTime.Hour() == 1 && nowTime.Minute() == 0 && checkWeekday(nowTime) {
		forceKenkou(session, DefaultGuildID, DefaultChannelID)
		return true
	} else {
		return false
	}
}

func forceKenkou(session *discordgo.Session, guildID string, channelID string) {
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

func playHotaru(session *discordgo.Session, guildID string, channelID string) {
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

func checkWeekday(time time.Time) bool {
	url := "https://s-proj.com/utils/checkHoliday.php?kind=h&date=" + time.Format("20060102")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray) == "else"
}
