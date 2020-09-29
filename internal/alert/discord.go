package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
	"github.com/ianmarmour/nvidia-clerk/internal/rest"
)

// DiscordMessage represents a discord message
type DiscordMessage interface {
	Get() string
	Set(string, string)
	JSON() ([]byte, error)
}

// DiscordAPIMessage represents a discord message relating to API status changes.
type DiscordAPIMessage struct {
	body string
}

// Get returns the current DiscordAPIMessage
func (d *DiscordAPIMessage) Get() string {
	return d.body
}

// Set takes in an API name and returns the JSON body for a Discord POST request
func (d *DiscordAPIMessage) Set(name string, status string) {
	d.body = fmt.Sprintf("NVIDIA API %s is now %s", name, status)
}

// JSON returns the JSON encoded bytes of a DiscordAPIMessage
func (d *DiscordAPIMessage) JSON() ([]byte, error) {
	body := map[string]string{"content": d.Get()}
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return json, nil
}

// DiscordProductMessage represents a discord message relating to product avaliablity.
type DiscordProductMessage struct {
	body string
}

// Get returns the current DiscordProductMessage
func (d *DiscordProductMessage) Get() string {
	return d.body
}

// Set takes in a product URL and returns the JSON body for a Discord POST request
func (d *DiscordProductMessage) Set(url string, status string) {
	d.body = "<" + url + ">"
}

// JSON returns the JSON encoded bytes of a DiscordProductMessage
func (d *DiscordProductMessage) JSON() ([]byte, error) {
	body := map[string]string{"content": d.Get()}
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return json, nil
}

//SendDiscordMessage Sends a notification message to a Discord Webhook.
func SendDiscordMessage(message DiscordMessage, config config.DiscordConfig, client *http.Client) error {
	json, err := message.JSON()
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

// StartDiscordAPINotifications Runs a loop and notifies discord when there is a status change.
func StartDiscordAPINotifications(region string, api string, config config.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{Timeout: 10 * time.Second}
	previousStatus := ""

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	check := make(chan bool)

	go func() {
		time.Sleep(61 * time.Second)
		check <- true
	}()

	for {
		select {
		case <-check:
			switch api {
			case "session":
				_, sessErr := rest.GetSessionToken(client)
				if sessErr != nil {
					if previousStatus != "offline" {
						message := DiscordAPIMessage{}
						message.Set("Store Session", "offline")
						SendDiscordMessage(&message, *config.DiscordConfig, client)
						previousStatus = "offline"
						log.Println(fmt.Sprintf("Sending Discord Notification for %s", config.Locale))
					}
				} else {
					if previousStatus != "online" {
						message := DiscordAPIMessage{}
						message.Set("Store Session", "online")
						SendDiscordMessage(&message, *config.DiscordConfig, client)
						previousStatus = "online"
						log.Println(fmt.Sprintf("Sending Discord Notification for %s", config.Locale))
					}
				}
			case "checkout":
				token, _ := rest.GetSessionToken(client)
				_, chkErr := rest.AddToCheckout(*config.SKU, token.Value, config.NvidiaLocale, client)
				if chkErr != nil {
					if previousStatus != "offline" {
						message := DiscordAPIMessage{}
						message.Set(fmt.Sprintf("%s Store Product Checkout", region), "offline")
						SendDiscordMessage(&message, *config.DiscordConfig, client)
						log.Println(fmt.Sprintf("Sending Discord Notification for %s", config.Locale))
						previousStatus = "offline"
					}
				} else {
					if previousStatus != "online" {
						message := DiscordAPIMessage{}
						message.Set(fmt.Sprintf("%s Store Product Checkout", region), "online")
						SendDiscordMessage(&message, *config.DiscordConfig, client)
						log.Println(fmt.Sprintf("Sending Discord Notification for %s", config.Locale))
						previousStatus = "online"
					}
				}
			}
		}
	}
}

// StartDiscordProductNotifications Runs a loop and notifies discord when there is a status change.
func StartDiscordProductNotifications(model string, config config.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{Timeout: 10 * time.Second}
	previousStatus := ""

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	check := make(chan bool)

	go func() {
		time.Sleep(61 * time.Second)
		check <- true
	}()

	for {
		select {
		case <-check:
			_, sessErr := rest.GetSessionToken(client)
			if sessErr != nil {
				info, err := rest.GetSkuInfo(*config.SKU, config.Locale, config.Currency, client)
				if err != nil {
					log.Println(fmt.Sprintf("Error attempting to get product information for %s in %s", *config.SKU, config.Locale))
					return
				}

				if info.Products.Product[0].InventoryStatus.Status == "PRODUCT_INVENTORY_IN_STOCK" {
					if previousStatus != "instock" {
						message := DiscordProductMessage{}
						message.Set(fmt.Sprintf("%s in stock now", model), "")
						SendDiscordMessage(&message, *config.DiscordConfig, client)
						previousStatus = "instock"
						log.Println(fmt.Sprintf("Sending Discord Notification for %s", config.Locale))
					}
				}
			}
		}
	}
}
