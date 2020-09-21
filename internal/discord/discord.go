package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendMessage Sends a notification message to a Discord Webhook.
func SendMessage(item string, config config.DiscordConfig, client *http.Client) error {
	content := fmt.Sprintf("%s Ready for Purchase", item)

	payload := map[string]interface{}{
		"content": &content,
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		log.Fatalln("Unable to marshal discord payload.")
	}

	req, _ := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(payloadJSON))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 400 {
		log.Fatal("Unable to send message to discord, bad request.")
	}

	return nil
}
