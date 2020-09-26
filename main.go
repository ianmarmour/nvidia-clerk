package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
	"github.com/ianmarmour/nvidia-clerk/internal/rest"
)

var testsHaveErrors bool

func main() {
	log.SetFlags(log.LstdFlags)

	var region string
	var model string
	var delay int64

	// Parse Argument Flags
	flag.StringVar(&region, "region", "USA", "3 Letter region code")
	flag.StringVar(&model, "model", "3080", "GPU Model number E.X. 3070, 3080, 3090")
	flag.Int64Var(&delay, "delay", 1, "Delay for refreshing in miliseconds")
	twitter := flag.Bool("twitter", false, "Enable Twitter Posts for whenever SKU is in stock.")
	twilio := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	discord := flag.Bool("discord", false, "Enable Discord webhook notifications for whenever SKU is in stock.")
	telegram := flag.Bool("telegram", false, "Enable Telegram webhook notifications for whenever SKU is in stock.")
	remote := flag.Bool("remote", false, "Enable notification only mode.")
	test := flag.Bool("test", false, "Enable remote mode for when you're away from computer.")
	flag.Parse()

	config, configErr := config.Get(region, model, delay, *twilio, *discord, *twitter, *telegram)
	if configErr != nil {
		log.Fatal(configErr)
	}
	client := &http.Client{Timeout: 10 * time.Second}

	token, err := rest.GetSessionToken(client)
	if err != nil {
		log.Println("Error getting session token from NVIDIA retrying...")
	}

	// For when NVIDIAs store APIs are down.
	for token == nil {
		sleep(delay)
		token, err = rest.GetSessionToken(client)
		if err != nil {
			log.Printf("Error getting session token from NVIDIA retrying...")
			continue
		}

		break
	}

	for {
		sleep(delay)

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

		if info.Products.Product[0].InventoryStatus.Status == "PRODUCT_INVENTORY_IN_STOCK" || *test == true {
			cart, err := rest.AddToCheckout(*config.SKU, token.Value, config.Locale, client)
			if err != nil {
				log.Println("Error adding card to checkout retrying...")
				continue
			}

			err = notify(info.Products.Product[0].Name, cart.URL, *remote, config, client)
			if err != nil {
				log.Println("Error attempting to send notification retrying...")
				continue
			}

			if *remote != true {
				err = openbrowser(cart.URL)
				if err != nil {
					log.Fatal("Error attempting to open browser.", err)
				}
			}

			break
		}
	}
}

func notify(id string, url string, remote bool, config *config.Config, client *http.Client) error {
	if remote != true {
		url = "Checkout avaliable on system running this program"
	}

	if config.TwilioConfig != nil {
		err := alert.SendText(id, url, *config.TwilioConfig, client)
		if err != nil {
			log.Println("Error sending SMS notification, retrying...")
			return err
		}
	}

	if config.TwitterConfig != nil {
		err := alert.SendTweet(id, url, *config.TwitterConfig)
		if err != nil {
			log.Println("Error sending Twitter notification, retrying...")
			return err
		}
	}

	if config.DiscordConfig != nil {
		err := alert.SendDiscordMessage(id, url, *config.DiscordConfig, client)
		if err != nil {
			log.Println("Error sending Discord notification, retrying...")
			return err
		}
	}

	if config.TelegramConfig != nil {
		err := alert.SendTelegramMessage(id, url, *config.TelegramConfig, client)
		if err != nil {
			log.Println("Error sending Telegram notification, retrying...")
			return err
		}
	}

	return nil
}

func openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil
}

func sleep(delay int64) {
	// Force a randomized jitter of up to 5 seconds to avoid looking like a bot.
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)

	ns := time.Duration(n) * time.Second
	ds := time.Duration(delay/1000) * time.Second
	time.Sleep(time.Duration(ns + ds))
}
