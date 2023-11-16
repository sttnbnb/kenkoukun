package kenkou

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/shmn7iii/kenkoukun/internal/db"
)

type KenkouSetting struct {
	GuildId   string
	ChannelId string
	// Time datetime
}

func UpdateGuildKenkouSetting(guildId string, channelId string) {
	dbc := db.Connect()
	defer dbc.Close()

	cmd := "INSERT OR REPLACE INTO kenkou_settings (guild_id, channel_id) VALUES (?, ?)"
	_, err := dbc.Exec(cmd, guildId, channelId)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetGuildKenkouSetting(guidlId string) (string, error) {
	dbc := db.Connect()
	defer dbc.Close()

	cmd := "SELECT * FROM kenkou_settings where guild_id = ?"
	row := dbc.QueryRow(cmd, guidlId)
	var setting KenkouSetting
	err := row.Scan(&setting.GuildId, &setting.ChannelId)
	if err != nil {
		return "", err
	}

	return setting.ChannelId, nil
}

func GetKenkouSettings() ([]KenkouSetting, error) {
	dbc := db.Connect()
	defer dbc.Close()

	cmd := "SELECT * FROM kenkou_settings"
	rows, _ := dbc.Query(cmd)
	defer rows.Close()
	var settings []KenkouSetting
	for rows.Next() {
		var setting KenkouSetting
		err := rows.Scan(&setting.GuildId, &setting.ChannelId)
		if err != nil {
			log.Println(err)
		}
		settings = append(settings, setting)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return settings, nil
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

func checkWeekday(time time.Time) bool {
	url := "https://s-proj.com/utils/checkHoliday.php?kind=h&date=" + time.Format("20060102")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray) == "else"
}
