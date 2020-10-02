package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/ianmarmour/nvidia-clerk/internal/alert"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
	"github.com/ianmarmour/nvidia-clerk/internal/rest"
	"github.com/ianmarmour/nvidia-clerk/internal/update"
)

func main() {
	log.SetFlags(log.LstdFlags)

	var region string
	var model string
	var delay int64

	// Parse Argument Flags
	flag.StringVar(&region, "region", "", "3 Letter region code E.X. USA, GBR, DEU")
	flag.StringVar(&model, "model", "", "GPU Model number E.X. 3070, 3080, 3090")
	flag.Int64Var(&delay, "delay", 1, "Delay for refreshing in miliseconds")
	twitter := flag.Bool("twitter", false, "Enable Twitter Posts for whenever SKU is in stock.")
	twilio := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	discord := flag.Bool("discord", false, "Enable Discord webhook notifications for whenever SKU is in stock.")
	telegram := flag.Bool("telegram", false, "Enable Telegram webhook notifications for whenever SKU is in stock.")
	remote := flag.Bool("remote", false, "Enable remote notification only mode.")
	desktop := flag.Bool("desktop", false, "Enable desktop notifications, disabled by default.")
	autoUpdate := flag.Bool("update", true, "Disable automatic updates, enabled by default.")
	test := flag.Bool("test", false, "Test notifications only, disabled by default")
	flag.Parse()

	config, configErr := config.Get(region, model, delay, *twilio, *discord, *twitter, *telegram, *desktop, false, *autoUpdate)
	if configErr != nil {
		log.Fatal(configErr)
	}
	client := &http.Client{Timeout: 10 * time.Second}

	if *test {
		var id string = fmt.Sprintf("TEST NOTIFICATION!!! %s", model)
		var url string = getURL(config.NvidiaLocale, model)
		notify(id, url, *remote, config, client)
	} else {
		var (
			mu    sync.Mutex
			token rest.SessionToken
		)
		var wg sync.WaitGroup
	
		wg.Add(3)
		go update.FetchApply(config.SystemConfig.UpdateURL, &wg)
		go getToken(client, delay, &token, &mu, &wg)
		go getGPU(client, config, model, *remote, delay, &wg)
	
		wg.Wait()
	}
	
}

func getToken(client *http.Client, delay int64, token *rest.SessionToken, mu *sync.Mutex, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		newToken, err := rest.GetSessionToken(client)
		if err != nil {
			log.Println("Error getting session token from NVIDIA retrying...")
		}

		if token == nil {
			mu.Lock()
			token = newToken
			mu.Unlock()
		} else {
			if token.Value != newToken.Value {
				mu.Lock()
				token = newToken
				mu.Unlock()
			}
		}

		sleep(delay)
	}
}

func getURL(locale string, model string) string {
	var cartURL string

	switch model {
	case "2060":
		cartURL = fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/rtx-%s-super/", locale, model)
	case "2070":
		cartURL = fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/rtx-%s-super/", locale, model)
	case "2080":
		cartURL = fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/rtx-%s-super/", locale, model)
	case "2080TI":
		cartURL = fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/rtx-%s-ti/", locale, model)
	case "3080":
		cartURL = fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/30-series/rtx-%s/", locale, model)
	case "3090":
		cartURL = fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/30-series/rtx-%s/", locale, model)
	default:
		cartURL = "https://www.nvidia.com/"
	}

	return cartURL
}

func getGPU(client *http.Client, config *config.Config, model string, remote bool, delay int64, wg *sync.WaitGroup) error {
	defer wg.Done()

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

		if info.Products.Product[0].InventoryStatus.Status == "PRODUCT_INVENTORY_IN_STOCK" {
			var cartURL string = getURL(config.NvidiaLocale, model);

			err = notify(info.Products.Product[0].Name, fmt.Sprintf(cartURL, model), remote, config, client)
			if err != nil {
				log.Println("Error attempting to send notification retrying...")
				continue
			}

			if remote != true {
				err = openbrowser(cartURL)
				if err != nil {
					log.Fatal("Error attempting to open browser.", err)
				}
			}

			break
		}
	}

	return nil
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
		message := alert.DiscordProductMessage{}
		message.Set(url, "IN_STOCK")

		err := alert.SendDiscordMessage(&message, *config.DiscordConfig, client)
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

	if config.ToastConfig != nil {
		err := alert.SendToast(config.ToastConfig.OS, id)
		if err != nil {
			log.Println("Error sending Windows Desktop notification, retrying...")
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
