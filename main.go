package main

import (
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
	BotToken  = os.Getenv("BOT_TOKEN")  //bot no token
	GuildID   = os.Getenv("GUILD_ID")   //kono guild desika ugokan w
	ChannelID = os.Getenv("CHANNEL_ID") //kono channel desika ugokan w
)

func main() {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return
	}

	session.AddHandler(SlashCommandsHandler)
	session.AddHandler(currentStatusNotification)

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	log.Println("Connection established.")
	log.Println("Hi there :)")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

loop:
	for {
		select {
		case <-sc:
			for _, v := range registeredCommands {
				err := session.ApplicationCommandDelete(session.State.User.ID, GuildID, v.ID)
				if err != nil {
					log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
				}
			}
			log.Println("interrupt. goodbye!")
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
	vc, _ := session.ChannelVoiceJoin(GuildID, ChannelID, false, true)
	go playHotaru(session, vc)
	log.Println("|_･) VC Joined.")
}

func forceKenkou(session *discordgo.Session) bool {
	if len(session.VoiceConnections) == 0 {
		return false
	}
	members, _ := session.GuildMembers(GuildID, "", 1000)
	for _, member := range members {
		session.GuildMemberMove(GuildID, member.User.ID, nil)
	}
	log.Println(">< All kicked.")
	return true
}

func playHotaru(session *discordgo.Session, vc *discordgo.VoiceConnection) {
	encodeSession, _ := dca.EncodeFile("./assets/hotaru.mp3", dca.StdEncodeOptions)
	vc.Speaking(true)
	done := make(chan error)
	dca.NewStream(encodeSession, vc, done)
	err := <-done
	if err != nil && err != io.EOF {
		log.Println("err", err)
	}
	vc.Speaking(false)
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "kenkou",
		Description: "Force Kenkou",
	},
	{
		Name:        "rename",
		Description: "Rename channel name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "role",
				Required:    true,
			},
		},
	},
	{
		Name:        "newname",
		Description: "Create new channel name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "New channel name",
				Required:    true,
			},
		},
	},
}

func SlashCommandsHandler(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "kenkou":
		go joinVC(session)
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

	case "rename":
		role := i.ApplicationCommandData().Options[0].RoleValue(session, i.GuildID)
		session.ChannelEdit(i.ChannelID, role.Name)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Channel name has changed to" + role.Name,
				Flags:   1 << 6,
			},
		})

	case "newname":
		name := i.ApplicationCommandData().Options[0].StringValue()
		session.ChannelEdit(i.ChannelID, name)

		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Channel name has changed to" + name,
				Flags:   1 << 6,
			},
		})
	}
}

func currentStatusNotification(session *discordgo.Session, m *discordgo.ChannelUpdate) {
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
	return
}
