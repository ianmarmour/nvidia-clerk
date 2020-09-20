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

		fmt.Println("Product Name: " + skuInfo.Products.Product[0].Name)
		fmt.Println("Product Status: " + skuInfo.Products.Product[0].InventoryStatus.Status)

		skuInventory, skuInventoryErr := rest.GetSkuInventory(config.SKU, httpClient)
		if skuInventoryErr != nil {
			fmt.Printf("Error getting SKU Inventory retrying...\n")
			continue
		}

		fmt.Printf("Product Inventory: %d", skuInventory.Product.AvailableQuantity)
		fmt.Printf("\n")
		fmt.Printf("\n")

		if skuInfo.Products.Product[0].InventoryStatus.Status == "PRODUCT_INVENTORY_IN_STOCK" {
			browser.NavigateTo(fmt.Sprintf("https://store.nvidia.com/store/nvidia/en_US/buy/productID.%s/clearCart.yes/nextPage.QuickBuyCartPage", config.SKU))

			textErr := alert.SendText(skuInfo.Products.Product[0].DisplayName, config.TwilioConfig, httpClient)
			if textErr != nil {
				fmt.Printf("Error sending notification retrying...\n")
				continue
			}

			os.Exit(0)
		}

		time.Sleep(time.Second)
	}
}
