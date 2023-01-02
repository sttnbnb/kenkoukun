package internal

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var HotaruDCABuffer = make([][]byte, 0)

func PlayHotaru(s *discordgo.Session, guildID, channelID string) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Fatalf("Can't get guild: %v", err)
		return
	}

	log.Println("|_･) VC Joined.")
	time.Sleep(250 * time.Millisecond)
	vc.Speaking(true)

	for _, buff := range HotaruDCABuffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)
	time.Sleep(250 * time.Millisecond)

	vc.Disconnect()
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
}

func CheckWeekday(time time.Time) bool {
	url := "https://s-proj.com/utils/checkHoliday.php?kind=h&date=" + time.Format("20060102")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray) == "else"
}

func LoadSound() error {
	file, err := os.Open("assets/hotaru.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		HotaruDCABuffer = append(HotaruDCABuffer, InBuf)
	}
}
