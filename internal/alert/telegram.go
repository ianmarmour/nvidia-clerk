package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

type TelegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

//SendTelegramMessage Sends a notification message to a Telegram Webhook.
func SendTelegramMessage(item string, config config.TelegramConfig, client *http.Client) error {
	content := fmt.Sprintf("%s Ready for Purchase", item)

	payload, err := json.Marshal(TelegramPayload{fmt.Sprintf("@%s", config.ChatID), content})
	if err != nil {
		log.Fatalln("Unable to marshal telegram payload.")
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.APIKey), bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Println("Unable to send message to telegram, bad request.")
		return errors.New("Bad Request")
	}

	return nil
}
