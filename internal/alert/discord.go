package alert

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendDiscordMessage Sends a notification message to a Discord Webhook.
func SendDiscordMessage(item string, nvidiaURL string, config config.DiscordConfig, client *http.Client) error {
	body := map[string]string{"content": item + " Ready for Purchase: " + nvidiaURL}

	json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	log.Println(json)

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode > 400 {
		return err
	}

	log.Println(r)

	defer r.Body.Close()

	return nil
}
