package main

import (
	"log"
	"sync"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
	"github.com/ianmarmour/nvidia-clerk/internal/rest"
)

func main() {
	var wg sync.WaitGroup

	cfg, err := config.Get("USA", "2060", 1, false, true, false, false, false, true)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(3)

	go alert.StartDiscordAPINotifications("session", *cfg, &wg)
	go alert.StartDiscordAPINotifications("checkout", *cfg, &wg)
	go rest.StartShieldsAPIServer(*cfg.ShieldsConfig, &wg)

	wg.Wait()
}
