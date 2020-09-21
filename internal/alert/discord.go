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

type DiscordPayload struct {
	Content string `json:"content"`
}

//SendDiscordMessage Sends a notification message to a Discord Webhook.
func SendDiscordMessage(item string, config config.DiscordConfig, client *http.Client) error {
	content := fmt.Sprintf("%s Ready for Purchase", item)

	payload, err := json.Marshal(DiscordPayload{content})

	if err != nil {
		log.Fatalln("Unable to marshal discord payload.")
	}

	req, _ := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 400 {
		fmt.Println("Unable to send message to discord, bad request.")
		return errors.New("Bad Request")
	}

	return nil
}
