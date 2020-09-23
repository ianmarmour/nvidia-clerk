package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

type payload struct {
	Content string `json:"content"`
}

//SendDiscordMessage Sends a notification message to a Discord Webhook.
func SendDiscordMessage(item string, config config.DiscordConfig, client *http.Client) error {
	body := fmt.Sprintf("%s Ready for Purchase", item)

	json, err := json.Marshal(payload{body})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(json))
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
