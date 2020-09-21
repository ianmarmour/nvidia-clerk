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

func main() {
	browser.StartChromeDebugMode()
	sessionContext := browser.StartSession()

	useSms := flag.Bool("sms", false, "Enable SMS notifications for whenever SKU is in stock.")
	flag.Parse()

	config := config.GetConfig(*useSms)

	httpClient := &http.Client{Timeout: 10 * time.Second}

	for {
		skuInfo, skuInfoErr := rest.GetSkuInfo(config.SKU, httpClient)
		if skuInfoErr != nil {
			fmt.Printf("Error getting SKU Information retrying...\n")
			continue
		}
		skuName := skuInfo.Products.Product[0].Name
		skuStatus := skuInfo.Products.Product[0].InventoryStatus.Status
		skuInventory, skuInventoryErr := rest.GetSkuInventory(config.SKU, httpClient)
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

			// Exit clean after a SKU was added to checkout cart.
			os.Exit(0)
		}

		time.Sleep(time.Second)
	}
}
