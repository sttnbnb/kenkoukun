package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

var (
	BotToken  = flag.String("token", "", "bot no token")
	GuildID   = flag.String("guild", "", "kono guild desika ugokan w")
	ChannelID = flag.String("channel", "", "kono channel desika ugokan w")
	RoleID    = flag.String("role", "", "kieru kikisen role")
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

	session.AddHandler(voiceStateUpdateEvent)
	session.AddHandler(SlashCommandsHandler)

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}

	for _, cmd := range commands {
		session.ApplicationCommandCreate(session.State.User.ID, *GuildID, cmd)
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
	vc, _ := session.ChannelVoiceJoin(*GuildID, *ChannelID, false, true)
	go playHotaru(session, vc)
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

func voiceStateUpdateEvent(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.ChannelID == *ChannelID {
		s.GuildMemberRoleAdd(*GuildID, v.UserID, *RoleID)
	} else if v.ChannelID == "" {
		s.GuildMemberRoleRemove(*GuildID, v.UserID, *RoleID)
	}
}

func playHotaru(session *discordgo.Session, vc *discordgo.VoiceConnection) {
	encodeSession, _ := dca.EncodeFile("./hotaru.mp3", dca.StdEncodeOptions)
	vc.Speaking(true)
	done := make(chan error)
	dca.NewStream(encodeSession, vc, done)
	err := <-done
	if err != nil && err != io.EOF {
		fmt.Print("err", err)
	}
	vc.Speaking(false)
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "kenkou",
		Description: "Force Kenkou",
	},
}

func SlashCommandsHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.ApplicationCommandData().Name
	if cmd == "kenkou" {
		joinVC(session)
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Let's Kenkou!",
				Flags:   1 << 6,
			},
		})
	loop:
		for {
			select {
			case <-time.After(5 * time.Minute):
				forceKenkou(session)
				break loop
			}
		}
	}
}
