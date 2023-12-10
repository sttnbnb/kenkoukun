package kenkou

import (
	"encoding/binary"
	"io"
	"net/http"
	"os"
	"time"
)

func LoadSound() error {
	file, err := os.Open("assets/hotaru.dca")
	if err != nil {
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
			return err
		}

		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		if err != nil {
			return err
		}

		HotaruDCABuffer = append(HotaruDCABuffer, InBuf)
	}
}

func checkWeekday(time time.Time) bool {
	url := "https://s-proj.com/utils/checkHoliday.php?kind=h&date=" + time.Format("20060102")
	resp, err := http.Get(url)
	if err != nil {
		// 通信失敗したら平日で返しちゃえ
		return true
	}

	defer resp.Body.Close()
	byteArray, _ := io.ReadAll(resp.Body)
	return string(byteArray) == "else"
}
