package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
	"github.com/ianmarmour/nvidia-clerk/internal/rest"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(len(config.RegionalConfigs) + 2)

	cfg, err := config.Get("USA", "2060", 1, false, true, false, false, false, true)
	if err != nil {
		log.Fatal(err)
	}

	// Monitor a single region for the session
	go alert.StartDiscordAPINotifications("USA", "session", *cfg, &wg)

	// Monitor only USA for the shields API sorry other regions.
	go rest.StartShieldsAPIServer(*cfg.ShieldsConfig, &wg)

	// Setup Notifications for all other regions avoiding rate limiting.
	for id := range config.RegionalConfigs {
		time.Sleep(10 * time.Second)
		tempID := id
		c, err := config.Get(tempID, "2060", 1, false, true, false, false, false, false)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting configuration for %s", tempID))
			wg.Add(-1)
			continue
		}

		dcURL, dcURLOK := os.LookupEnv(fmt.Sprintf("DISCORD_WEBHOOK_URL_%s", tempID))
		if dcURLOK == false {
			log.Println(fmt.Sprintf("Error getting discord webhook configuration for %s", tempID))
			wg.Add(-1)
			continue
		}

		cc := *c
		cc.DiscordConfig = &config.DiscordConfig{WebhookURL: dcURL}

		log.Println(fmt.Sprintf("Starting goroutine for %s", tempID))
		go alert.StartDiscordAPINotifications(tempID, "checkout", cc, &wg)
	}

	wg.Add(len(config.RegionalConfigs))
	// Setup 3080 Product Notifications for all regions
	for id := range config.RegionalConfigs {
		time.Sleep(10 * time.Second)
		tempID := id
		c, err := config.Get(tempID, "3080", 1, false, true, false, false, false, false)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting configuration for %s", tempID))
			wg.Add(-1)
			continue
		}

		dcURL, dcURLOK := os.LookupEnv(fmt.Sprintf("DISCORD_WEBHOOK_URL_%s", tempID))
		if dcURLOK == false {
			log.Println(fmt.Sprintf("Error getting discord webhook configuration for %s", tempID))
			wg.Add(-1)
			continue
		}

		cc := *c
		cc.DiscordConfig = &config.DiscordConfig{WebhookURL: dcURL}

		log.Println(fmt.Sprintf("Starting goroutine for %s", tempID))
		go alert.StartDiscordProductNotifications("3080", cc, &wg)
	}

	wg.Add(len(config.RegionalConfigs))
	// Setup 3090 Product Notifications for all regions
	for id := range config.RegionalConfigs {
		time.Sleep(10 * time.Second)
		tempID := id
		c, err := config.Get(tempID, "3090", 1, false, true, false, false, false, false)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting configuration for %s", tempID))
			wg.Add(-1)
			continue
		}

		dcURL, dcURLOK := os.LookupEnv(fmt.Sprintf("DISCORD_WEBHOOK_URL_%s", tempID))
		if dcURLOK == false {
			log.Println(fmt.Sprintf("Error getting discord webhook configuration for %s", tempID))
			wg.Add(-1)
			continue
		}

		cc := *c
		cc.DiscordConfig = &config.DiscordConfig{WebhookURL: dcURL}

		log.Println(fmt.Sprintf("Starting goroutine for %s", tempID))
		go alert.StartDiscordProductNotifications("3090", cc, &wg)
	}

	wg.Wait()
}
