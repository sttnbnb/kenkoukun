package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/shmn7iii/kenkoukun/internal"
	"github.com/shmn7iii/kenkoukun/internal/chatgpt"
	"github.com/shmn7iii/kenkoukun/internal/kenkou"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken         = os.Getenv("BOT_TOKEN")          //bot no token
	DefaultGuildID   = os.Getenv("DEFAULT_GUILD_ID")   //teizi kenkou guild
	DefaultChannelID = os.Getenv("DEFAULT_CHANNEL_ID") //teizi kenkou channel
)

func main() {
	var err error

	err = kenkou.LoadSound()
	if err != nil {
		log.Fatalf("Error loading sound: %v", err)
		return
	}

	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return
	}

	session.AddHandler(internal.SlashCommandHandler)
	session.AddHandler(chatgpt.MessageCreateEventHandler)

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(internal.Commands))
	for i, v := range internal.Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	log.Println("Hi there :)")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)

loop:
	for {
		select {
		case <-sc:
			stop(session, registeredCommands)
			break loop
		case <-time.After(59 * time.Second):
			kenkou.KenkouBatch(session, DefaultGuildID, DefaultChannelID)
		}
	}
}

func stop(session *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand) {
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("interrupt. goodbye!")
	session.Close()
}
