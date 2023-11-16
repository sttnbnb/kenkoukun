package kenkou

import (
	"log"

	"github.com/shmn7iii/kenkoukun/internal/db"
)

type KenkouSetting struct {
	GuildId   string
	ChannelId string
}

func UpdateGuildKenkouSetting(setting KenkouSetting) {
	dbc := db.Connect()
	defer dbc.Close()

	cmd := "INSERT OR REPLACE INTO kenkou_settings (guild_id, channel_id) VALUES (?, ?)"
	_, err := dbc.Exec(cmd, setting.GuildId, setting.ChannelId)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetGuildKenkouSetting(guidlId string) (KenkouSetting, error) {
	dbc := db.Connect()
	defer dbc.Close()

	cmd := "SELECT * FROM kenkou_settings where guild_id = ?"
	row := dbc.QueryRow(cmd, guidlId)
	var setting KenkouSetting
	err := row.Scan(&setting.GuildId, &setting.ChannelId)
	if err != nil {
		return KenkouSetting{}, err
	}

	return setting, nil
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
