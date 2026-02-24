package fnNotifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendTelegramMessage(token, chatID, message string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[Telegram] 전송 실패 (Network): %v\n", err)
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("[Telegram] 응답 바디 닫기 실패: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Telegram] 전송 실패 (Status %d)\n", resp.StatusCode)
	}
}
