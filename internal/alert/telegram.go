package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendTelegramMessage Sends a notification message to a Telegram Webhook.
func SendTelegramMessage(item string, nvidiaURL string, config config.TelegramConfig, client *http.Client) error {
	text := item + " Ready for Purchase: " + nvidiaURL
	body := map[string]interface{}{"chat_id": config.ChatID, "text": text, "disable_web_page_preview": true}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// We're required to disable web page previews to ensure that the cart links don't get invalidated
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.APIKey), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return nil
}
