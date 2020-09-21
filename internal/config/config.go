package config

import (
	"fmt"
	"log"
	"os"
)

type RegionalConfig struct {
	SKU      string
	Locale   string
	Currency string
}

type TwilioConfig struct {
	AccountSID        string
	Token             string
	SourceNumber      string
	DestinationNumber string
}

type Config struct {
	Locale       string
	Currency     string
	SKU          string
	TwilioConfig TwilioConfig
}

var skuBasedConfig = map[string]RegionalConfig{
	"5438481700": {
		SKU:      "5438481700",
		Locale:   "en_us",
		Currency: "USD",
	},
	"5438792800": {
		SKU:      "5438792800",
		Locale:   "en_gb",
		Currency: "GBP",
	},
	"5438792300": {
		SKU:      "5438792300",
		Locale:   "de_de",
		Currency: "EUR",
	},
	"5438795200": {
		SKU:      "5438795200",
		Locale:   "fr_fr",
		Currency: "EUR",
	},
}

//GetTwilioConfig Generates TwilioConfiguration for application from environmental variables.
func GetTwilioConfig() TwilioConfig {
	configuration := TwilioConfig{}

	accountSid, accountSidOk := os.LookupEnv("TWILIO_ACCOUNT_SID")
	if accountSidOk == false {
		log.Fatal("TWLIO_ACCOUNT_SID Environment Variable is unset, exiting.")
	}

	configuration.AccountSID = accountSid

	token, tokenOk := os.LookupEnv("TWILIO_TOKEN")
	if tokenOk == false {
		log.Fatal("TWILIO_TOKEN Environment Variable is unset, exiting.")
	}

	configuration.Token = token

	sourceNumber, sourceNumberOk := os.LookupEnv("TWILIO_SOURCE_NUMBER")
	if sourceNumberOk == false {
		log.Fatal("TWILIO_SOURCE_NUMBER Environment Variable is unset, exiting.")
	}

	configuration.SourceNumber = sourceNumber

	destinationNumber, destinationNumberOk := os.LookupEnv("TWILIO_DESTINATION_NUMBER")
	if destinationNumberOk == false {
		log.Fatal("TWILIO_DESTINATION_NUMBER Environment Variable is unset, exiting.")
	}

	configuration.DestinationNumber = destinationNumber

	return configuration
}

//GetConfig Generates Configuration for application from environmental variables.
func GetConfig(smsEnabled bool) Config {
	configuration := Config{}

	sku, skuOk := os.LookupEnv("NVIDIA_CLERK_SKU")
	if skuOk == false {
		log.Fatal("NVIDIA_CLERK_SKU Environment Variable is unset, exiting.")
	}

	configuration.SKU = sku

	locale, localeOk := os.LookupEnv("NVIDIA_CLERK_LOCALE")
	if localeOk == false {
		locale = skuBasedConfig[sku].Locale
		fmt.Println(fmt.Sprintf("NVIDIA_CLERK_LOCALE unset defaulting locale to %s based on SKU", locale))
	}

	configuration.Locale = locale

	currency, currencyOk := os.LookupEnv("NVIDIA_CLERK_CURRENCY")
	if currencyOk == false {
		currency = skuBasedConfig[sku].Currency
		fmt.Println(fmt.Sprintf("NVIDIA_CLERK_CURRENCY unset defaulting currency to %s based on SKU", currency))
	}

	configuration.Currency = currency

	if smsEnabled == true {
		configuration.TwilioConfig = GetTwilioConfig()
	}

	return configuration
}
