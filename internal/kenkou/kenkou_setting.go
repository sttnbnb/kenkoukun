package kenkou

import (
	"time"

	"github.com/shmn7iii/kenkoukun/internal/database"
	"gorm.io/gorm"
)

type KenkouSetting struct {
	gorm.Model
	GuildId   string `gorm:"primaryKey"`
	ChannelId *string
	Time      time.Time `gorm:"default:0000-01-01 01:00:00"` // TimeOnly
}

func init() {
	err := database.Db.AutoMigrate(&KenkouSetting{})
	if err != nil {
		panic("Failed to migrate database.")
	}
}

func SaveGuildKenkouSetting(newSetting KenkouSetting) {
	var setting KenkouSetting
	// UPSERT
	database.Db.Where(KenkouSetting{GuildId: newSetting.GuildId}).Assign(newSetting).FirstOrCreate(&setting)
}

func GetGuildKenkouSetting(guildId string) (KenkouSetting, error) {
	var setting KenkouSetting
	database.Db.FirstOrCreate(&setting, KenkouSetting{GuildId: guildId})
	return setting, nil
}

func GetKenkouSettings() ([]KenkouSetting, error) {
	var settings []KenkouSetting
	database.Db.Find(&settings)
	return settings, nil
}
