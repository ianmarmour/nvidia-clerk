package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

type TelegramPayload struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

//SendTelegramMessage Sends a notification message to a Telegram Webhook.
func SendTelegramMessage(item string, config config.TelegramConfig, client *http.Client) error {
	body := fmt.Sprintf("%s Ready for Purchase", item)
	payload, err := json.Marshal(TelegramPayload{config.ChatID, body})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.APIKey), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
