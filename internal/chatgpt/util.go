package chatgpt

import (
	"log"
	"os"
)

var (
	chatSystemPrompt string
)

func init() {
	loadChatSystemPrompt()
}

// 外部に用意した設定ファイルを読み込む
func loadChatSystemPrompt() {
	f, err := os.Open("assets/chat_system_prompt.txt")
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}

	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		if err.Error() == "EOF" {
			chatSystemPrompt = ""
			return
		}
		log.Fatalf("Cannot read file: %v", err)
	}

	chatSystemPrompt = string(data[:count])
}
