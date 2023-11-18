package kenkou

import (
	"time"

	"github.com/shmn7iii/kenkoukun/internal/database"
	"gorm.io/gorm"
)

// ChannelId が nil なのに KenkouAlarm が true になり得るのがキモい
type KenkouSetting struct {
	gorm.Model
	Guild            string `gorm:"primaryKey"`
	AlarmActive      bool   `gorm:"default:true"`
	AlarmChannel     *string
	AlarmTime        time.Time `gorm:"default:0000-01-01 01:00:00"` // TimeOnly
	AlarmWeekdayOnly bool      `gorm:"default:true"`
}

func init() {
	err := database.Db.AutoMigrate(&KenkouSetting{})
	if err != nil {
		panic("Failed to migrate database.")
	}
}

func SaveGuildKenkouSetting(newSetting KenkouSetting) KenkouSetting {
	database.Db.Save(&newSetting)

	return newSetting
}

func GetGuildKenkouSetting(guildId string) (KenkouSetting, error) {
	var setting KenkouSetting
	database.Db.FirstOrCreate(&setting, KenkouSetting{Guild: guildId})
	return setting, nil
}

func GetKenkouSettings() ([]KenkouSetting, error) {
	var settings []KenkouSetting
	database.Db.Find(&settings)
	return settings, nil
}
