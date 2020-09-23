package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/browser"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

var testsHaveErrors bool

func runTest(name string, client *http.Client, config config.Config) {
	switch name {
	case "twilio":
		textErr := alert.SendText(config.SKU, config.TwilioConfig, client)
		if textErr != nil {
			fmt.Printf("Error testing SMS notification...\n")
			testsHaveErrors = true
		} else {
			fmt.Printf("SMS Notification testing completed successfully\n")
		}
	case "discord":
		discordErr := alert.SendDiscordMessage(config.SKU, config.DiscordConfig, client)
		if discordErr != nil {
			fmt.Printf("Error testing Discord notification...\n")
			testsHaveErrors = true
		} else {
			fmt.Printf("Discord Notification testing completed successfully\n")
		}
	case "twitter":
		tweetErr := alert.SendTweet(config.SKU, config.TwitterConfig)
		if tweetErr != nil {
			fmt.Printf("Error testing Twitter notification exiting...\n")
			os.Exit(1)
		} else {
			fmt.Printf("Twitter Notification testing completed succesfully")
		}
	case "telegram":
		telegramErr := alert.SendTelegramMessage(config.SKU, config.TelegramConfig, client)
		if telegramErr != nil {
			fmt.Printf("Error testing Telegram notification exiting...\n")
			os.Exit(1)
		} else {
			fmt.Printf("Telegram Notification testing completed succesfully")
		}
	default:

	}
}

func main() {
	log.SetFlags(log.LstdFlags)

	var region string
	var delay int64

	// Parse Argument Flags
	flag.StringVar(&region, "region", "USA", "3 Letter region code")
	flag.Int64Var(&delay, "delay", 500, "Delay for refreshing in miliseconds")
	twitter := flag.Bool("twitter", false, "Enable Twitter Posts for whenever SKU is in stock.")
	twilio := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	discord := flag.Bool("discord", false, "Enable Discord webhook notifications for whenever SKU is in stock.")
	telegram := flag.Bool("telegram", false, "Enable Telegram webhook notifications for whenever SKU is in stock.")
	test := flag.Bool("test", false, "Enable testing mode")
	flag.Parse()

	config, configErr := config.Get(region, delay, *twilio, *discord, *twitter, *telegram)
	if configErr != nil {
		log.Fatal(configErr)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Execute Tests
	if *test == true {
		config.SKU = config.TestSKU
		if *twilio == true {
			runTest("twilio", client, *config)
		}

		if *discord == true {
			runTest("discord", client, *config)
		}

		if *twitter == true {
			runTest("twitter", client, *config)
		}

		if *telegram == true {
			runTest("telegram", client, *config)
		}

		if testsHaveErrors == true {
			log.Printf("Testing failed with errors, exiting...\n")
			os.Exit(1)
		}
	}

	ctx, err := browser.Start(*config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		invStatus, invStatusErr := browser.GetInventoryStatus(ctx, config.SKU, config.Locale, config.Delay)
		if invStatusErr != nil {
			log.Printf("Error getting SKU Information retrying...\n")
			continue
		}
		id := invStatus.Product.ID
		status := invStatus.Status

		log.Println("Product ID: " + id)
		log.Println("Product Status: " + status)

		if status == "PRODUCT_INVENTORY_IN_STOCK" {
			cartErr := browser.AddToCart(ctx, config.SKU, config.Locale)

			if cartErr != nil {
				log.Printf("Error adding item to cart retrying...\n")
				continue
			}

			checkoutErr := browser.Checkout(ctx, config.Locale)
			if checkoutErr != nil {
				log.Printf("Error adding item to checkout retrying...\n")
				continue
			}

			if *twilio == true {
				textErr := alert.SendText(id, config.TwilioConfig, client)
				if textErr != nil {
					log.Printf("Error sending SMS notification retrying...\n")
					continue
				}
			}

			if *twitter == true {
				tweetErr := alert.SendTweet(id, config.TwitterConfig)
				if tweetErr != nil {
					log.Printf("Error sending Twitter notification retrying...\n")
					continue
				}
			}

			if *discord == true {
				discordErr := alert.SendDiscordMessage(id, config.DiscordConfig, client)
				if discordErr != nil {
					log.Printf("Error sending discord notification retrying...\n")
					continue
				}
			}

			if *telegram == true {
				telegramErr := alert.SendTelegramMessage(id, config.TelegramConfig, client)
				if telegramErr != nil {
					log.Printf("Error sending telegram notification retrying...\n")
					continue
				}
			}

			break
		}
	}

	os.Exit(0)
}
