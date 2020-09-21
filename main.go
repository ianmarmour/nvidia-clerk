package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/browser"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
	"github.com/ianmarmour/nvidia-clerk/internal/rest"
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
	default:

	}
}

func main() {
	// Parse Argument Flags
	useSms := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	useDiscord := flag.Bool("discord", false, "Enable Discord webhook notifications for whenever SKU is in stock.")
	useTest := flag.Bool("test", false, "Enable testing mode")
	flag.Parse()

	config := config.GetConfig(*useSms, *useDiscord)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	// Execute Tests
	if *useTest == true {
		if *useSms == true {
			runTest("sms", httpClient, config)
		}

		if *useDiscord == true {
			runTest("discord", httpClient, config)
		}

		if testsHaveErrors == true {
			fmt.Printf("Testing failed with errors, exiting...\n")
			os.Exit(1)
		} else {
			fmt.Printf("Testing completed succesfully, exiting...\n")
			os.Exit(0)
		}
	}

	browser.StartChromeDebugMode()
	sessionContext := browser.StartSession(config)

	for {
		skuInfo, skuInfoErr := rest.GetSkuInfo(config.SKU, config.Locale, config.Currency, httpClient)
		if skuInfoErr != nil {
			fmt.Printf("Error getting SKU Information retrying...\n")
			continue
		}
		skuName := skuInfo.Products.Product[0].Name
		skuStatus := skuInfo.Products.Product[0].InventoryStatus.Status
		skuInventory, skuInventoryErr := rest.GetSkuInventory(config.SKU, config.Locale, httpClient)
		if skuInventoryErr != nil {
			fmt.Printf("Error getting SKU Inventory retrying...\n")
			continue
		}

		fmt.Println("SKU Name: " + skuName)
		fmt.Println("SKU Status: " + skuStatus)
		fmt.Printf("Product Inventory: %d \n\n", skuInventory.Product.AvailableQuantity)

		if skuStatus == "PRODUCT_INVENTORY_IN_STOCK" {
			browser.AddToCart(sessionContext, config.SKU)
			browser.Checkout(sessionContext)

			if *useSms == true {
				textErr := alert.SendText(skuName, config.TwilioConfig, httpClient)
				if textErr != nil {
					fmt.Printf("Error sending notification retrying...\n")
					continue
				}
			}

			if *useDiscord == true {
				textErr := alert.SendDiscordMessage(skuName, config.DiscordConfig, httpClient)
				if textErr != nil {
					fmt.Printf("Error sending discord notification retrying...\n")
					continue
				}
			}

			// Exit clean after a SKU was added to checkout cart.
			os.Exit(0)
		}

		time.Sleep(time.Second)
	}
}
