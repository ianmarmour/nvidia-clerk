package alert

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendDiscordMessage Sends a notification message to a Discord Webhook.
func SendDiscordMessage(item string, nvidiaURL string, config config.DiscordConfig, client *http.Client) error {
	// The wrapping tags for the URL here stop discord from pulling metadata from our link disabling the checkout functionality.
	base := "<" + nvidiaURL + ">"
	body := map[string]string{"content": base}

	json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.WebhookURL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode > 400 {
		return err
	}

	defer r.Body.Close()

	return nil
}
