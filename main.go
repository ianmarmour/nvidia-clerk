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
	case "sms":
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
	default:

	}
}

func main() {
	var region string

	// Parse Argument Flags
	flag.StringVar(&region, "region", "USA", "3 Letter region code")
	useTwitter := flag.Bool("twitter", false, "Enable Twitter Posts for whenever SKU is in stock.")
	useSms := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	useDiscord := flag.Bool("discord", false, "Enable Discord webhook notifications for whenever SKU is in stock.")
	useTest := flag.Bool("test", false, "Enable testing mode")
	flag.Parse()

	config, configErr := config.GetConfig(region, *useSms, *useDiscord, *useTwitter)
	if configErr != nil {
		log.Fatal(configErr)
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	// Execute Tests
	if *useTest == true {
		config.SKU = config.TestSKU
		if *useSms == true {
			runTest("sms", httpClient, *config)
		}

		if *useDiscord == true {
			runTest("discord", httpClient, *config)
		}

		if *useTwitter == true {
			runTest("twitter", httpClient, *config)
		}

		if testsHaveErrors == true {
			fmt.Printf("Testing failed with errors, exiting...\n")
			os.Exit(1)
		}
	}

	sessionContext := browser.StartSession(*config)

	for {
		skuInfo, skuInfoErr := browser.GetInventoryStatus(sessionContext, config.SKU, config.Locale)
		if skuInfoErr != nil {
			fmt.Printf("Error getting SKU Information retrying...\n")
			continue
		}
		productID := skuInfo.Product.ID
		skuStatus := skuInfo.Status

		fmt.Println("Product ID: " + productID)
		fmt.Println("Product Status: " + skuStatus)

		if skuStatus == "PRODUCT_INVENTORY_IN_STOCK" {
			browser.AddToCart(sessionContext, config.SKU, config.Locale)
			browser.Checkout(sessionContext, config.Locale)

			if *useSms == true {
				textErr := alert.SendText(productID, config.TwilioConfig, httpClient)
				if textErr != nil {
					fmt.Printf("Error sending SMS notification retrying...\n")
					continue
				}
			}

			if *useTwitter == true {
				tweetErr := alert.SendTweet(productID, config.TwitterConfig)
				if tweetErr != nil {
					fmt.Printf("Error sending Twitter notification retrying...\n")
					continue
				}
			}

			if *useDiscord == true {
				discordErr := alert.SendDiscordMessage(productID, config.DiscordConfig, httpClient)
				if discordErr != nil {
					fmt.Printf("Error sending discord notification retrying...\n")
					continue
				}
			}

			// Exit clean after a SKU was added to checkout cart.
			os.Exit(0)
		}
	}
}
