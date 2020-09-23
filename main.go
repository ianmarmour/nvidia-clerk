package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/browser"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

var testsHaveErrors bool

func runTest(name string, client *http.Client, config config.Config) {
	switch name {
	case "twilio":
		textErr := alert.SendText(*config.SKU, config.TwilioConfig, client)
		if textErr != nil {
			fmt.Printf("Error testing SMS notification...\n")
			testsHaveErrors = true
		} else {
			fmt.Printf("SMS Notification testing completed successfully\n")
		}
	case "discord":
		discordErr := alert.SendDiscordMessage(*config.SKU, config.DiscordConfig, client)
		if discordErr != nil {
			fmt.Printf("Error testing Discord notification...\n")
			testsHaveErrors = true
		} else {
			fmt.Printf("Discord Notification testing completed successfully\n")
		}
	case "twitter":
		tweetErr := alert.SendTweet(*config.SKU, config.TwitterConfig)
		if tweetErr != nil {
			fmt.Printf("Error testing Twitter notification exiting...\n")
			os.Exit(1)
		} else {
			fmt.Printf("Twitter Notification testing completed succesfully")
		}
	case "telegram":
		telegramErr := alert.SendTelegramMessage(*config.SKU, config.TelegramConfig, client)
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
	var model string
	var delay int64

	// Parse Argument Flags
	flag.StringVar(&region, "region", "USA", "3 Letter region code")
	flag.StringVar(&model, "model", "3080", "GPU Model number E.X. 3070, 3080, 3090")
	flag.Int64Var(&delay, "delay", 500, "Delay for refreshing in miliseconds")
	twitter := flag.Bool("twitter", false, "Enable Twitter Posts for whenever SKU is in stock.")
	twilio := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	discord := flag.Bool("discord", false, "Enable Discord webhook notifications for whenever SKU is in stock.")
	telegram := flag.Bool("telegram", false, "Enable Telegram webhook notifications for whenever SKU is in stock.")
	test := flag.Bool("test", false, "Enable testing mode")
	flag.Parse()

	config, configErr := config.Get(region, model, delay, *twilio, *discord, *twitter, *telegram)
	if configErr != nil {
		log.Fatal(configErr)
	}

	// Execute Tests
	if *test == true {
		ExecuteTests(config, *twilio, *discord, *twitter, *telegram)
	}

	ctx, err := browser.Start(*config)
	if err != nil {
		log.Fatal(err)
	}

	if config.SKU == nil {
		LookupSKU(ctx, model, config)
	}

	MonitorRelease(ctx, config, *twilio, *discord, *twitter, *telegram)

	os.Exit(0)
}

func ExecuteTests(config *config.Config, twilio bool, discord bool, twitter bool, telegram bool) {
	client := &http.Client{Timeout: 10 * time.Second}

	config.SKU = config.TestSKU
	if twilio == true {
		runTest("twilio", client, *config)
	}

	if discord == true {
		runTest("discord", client, *config)
	}

	if twitter == true {
		runTest("twitter", client, *config)
	}

	if telegram == true {
		runTest("telegram", client, *config)
	}

	if testsHaveErrors == true {
		log.Printf("Testing failed with errors, exiting...\n")
		os.Exit(1)
	}
}

// LookupSKU Looks up the SKU based on model name of a GPU in Digital River
func LookupSKU(ctx context.Context, model string, config *config.Config) error {
	for {
		products, productsErr := browser.GetProducts(ctx, config.Locale, config.Delay)
		if productsErr != nil {
			log.Printf("Error getting Products Information retrying...\n")
		}

		for _, product := range products {
			if strings.Contains(product.DisplayName, model) == true {
				subs := strings.Split(product.URI, "/")
				config.SKU = &subs[len(subs)-1]
				return nil
			}
		}
	}
}

// MonitorRelease Monitors for GPU launches.
func MonitorRelease(ctx context.Context, config *config.Config, twilio bool, discord bool, twitter bool, telegram bool) {
	client := &http.Client{Timeout: 10 * time.Second}

	product, productErr := browser.GetProduct(ctx, *config.SKU, config.Locale, config.Delay)
	if productErr != nil {
		log.Printf("Error getting Product Information retrying...\n")
	}

	for {
		invStatus, invStatusErr := browser.GetInventoryStatus(ctx, *config.SKU, config.Locale, config.Delay)
		if invStatusErr != nil {
			log.Printf("Error getting Iventory Status retrying...\n")
			continue
		}
		id := invStatus.Product.ID
		status := invStatus.Status

		log.Println("Product ID: " + id)
		log.Println("Product Name: " + product.Name)
		log.Println("Product Locale: " + config.Locale)
		log.Println("Product Status: " + status + "\n")

		if status == "PRODUCT_INVENTORY_IN_STOCK" {
			cartErr := browser.AddToCart(ctx, *config.SKU, config.Locale)

			if cartErr != nil {
				log.Printf("Error adding item to cart retrying...\n")
				continue
			}

			checkoutErr := browser.Checkout(ctx, config.Locale)
			if checkoutErr != nil {
				log.Printf("Error adding item to checkout retrying...\n")
				continue
			}

			if twilio == true {
				textErr := alert.SendText(id, config.TwilioConfig, client)
				if textErr != nil {
					log.Printf("Error sending SMS notification retrying...\n")
					continue
				}
			}

			if twitter == true {
				tweetErr := alert.SendTweet(id, config.TwitterConfig)
				if tweetErr != nil {
					log.Printf("Error sending Twitter notification retrying...\n")
					continue
				}
			}

			if discord == true {
				discordErr := alert.SendDiscordMessage(id, config.DiscordConfig, client)
				if discordErr != nil {
					log.Printf("Error sending discord notification retrying...\n")
					continue
				}
			}

			if telegram == true {
				telegramErr := alert.SendTelegramMessage(id, config.TelegramConfig, client)
				if telegramErr != nil {
					log.Printf("Error sending telegram notification retrying...\n")
					continue
				}
			}

			break
		}
	}
}
