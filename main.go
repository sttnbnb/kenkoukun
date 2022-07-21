package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken  = flag.String("token", "", "bot no token")
	GuildID   = flag.String("guild", "", "kono guild desika ugokan w")
	ChannelID = flag.String("channel", "", "kono channel desika ugokan w")
)

func init() {
	flag.Parse()
}

func main() {
	session, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return
	}

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}

	fmt.Println("Connection established.")
	fmt.Println("Hi there :)")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

loop:
	for {
		select {
		case <-sc:
			fmt.Println("interrupt. goodbye!")
			session.Close()
			break loop
		case <-time.After(59 * time.Second):
			run(session)
		}
	}
}

func run(session *discordgo.Session) bool {
	nowTime := time.Now()
	if nowTime.Hour() == 0 && nowTime.Minute() == 55 && checkWeekday(nowTime) {
		joinVC(session)
		return true
	} else if nowTime.Hour() == 1 && nowTime.Minute() == 0 && checkWeekday(nowTime) {
		forceKenkou(session)
		return true
	} else {
		return false
	}
}

func checkWeekday(time time.Time) bool {
	url := "https://s-proj.com/utils/checkHoliday.php?kind=h&date=" + time.Format("20060102")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray) == "else"
}

func joinVC(session *discordgo.Session) {
	session.ChannelVoiceJoin(*GuildID, *ChannelID, true, true)
	fmt.Println("|_･) VC Joined.")
}

func forceKenkou(session *discordgo.Session) bool {
	if len(session.VoiceConnections) == 0 {
		return false
	}
	members, _ := session.GuildMembers(*GuildID, "", 1000)
	for _, member := range members {
		session.GuildMemberMove(*GuildID, member.User.ID, nil)
	}
	fmt.Println(">< All kicked.")
	return true
}
