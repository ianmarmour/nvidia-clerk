package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
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

	client := &http.Client{Timeout: 10 * time.Second}

	for {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(5)
		time.Sleep(time.Duration(n) * time.Second)

		info, err := rest.GetSkuInfo(*config.SKU, config.Locale, config.Currency, client)
		if err != nil {
			continue
		}

		// HACK: Resolves https://github.com/ianmarmour/nvidia-clerk/issues/85
		if len(info.Products.Product) < 1 {
			log.Printf("Error attempting to get product information retrying...\n")
			continue
		}

		log.Println(fmt.Sprintf("Product ID: %v", info.Products.Product[0].ID))
		log.Println("Product Name: " + info.Products.Product[0].Name)
		log.Println("Product Locale: " + config.Locale)
		log.Println("Product Status: " + info.Products.Product[0].InventoryStatus.Status + "\n")

		if info.Products.Product[0].InventoryStatus.Status == "PRODUCT_INVENTORY_IN_STOCK" {
			id := info.Products.Product[0].Name

			if *twilio == true {
				textErr := alert.SendText(id, config.TwilioConfig, client)
				if textErr != nil {
					log.Printf("Error sending SMS notification retrying...\n")
				}
			}

			if *twitter == true {
				tweetErr := alert.SendTweet(id, config.TwitterConfig)
				if tweetErr != nil {
					log.Printf("Error sending Twitter notification retrying...\n")
				}
			}

			if *discord == true {
				discordErr := alert.SendDiscordMessage(id, config.DiscordConfig, client)
				if discordErr != nil {
					log.Printf("Error sending discord notification retrying...\n")
				}
			}

			if *telegram == true {
				telegramErr := alert.SendTelegramMessage(id, config.TelegramConfig, client)
				if telegramErr != nil {
					log.Printf("Error sending telegram notification retrying...\n")
				}
			}

			ctx, err := browser.Start(*config)
			if err != nil {
				log.Fatal("Error attempting to open browser.")
			}
			browser.OpenProductPage(ctx, model, config.NvidiaLocale, *test)

			break
		}
	}
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
