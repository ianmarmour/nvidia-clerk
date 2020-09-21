package config

import (
	"fmt"
	"log"
	"os"
)

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

	locale, localeOk := os.LookupEnv("NVIDIA_CLERK_LOCALE")
	if localeOk == false {
		locale = "en_us"
		fmt.Println("NVIDIA_CLERK_LOCALE unset defaulting locale to en_us.")
	}

	configuration.Locale = locale

	currency, currencyOk := os.LookupEnv("NVIDIA_CLERK_CURRENCY")
	if currencyOk == false {
		currency = "USD"
		fmt.Println("NVIDIA_CLERK_CURRENCY unset defaulting currency to USD.")
	}

	configuration.Currency = currency

	sku, skuOk := os.LookupEnv("NVIDIA_CLERK_SKU")
	if skuOk == false {
		log.Fatal("NVIDIA_CLERK_SKU Environment Variable is unset, exiting.")
	}

	configuration.SKU = sku

	if smsEnabled == true {
		configuration.TwilioConfig = GetTwilioConfig()
	}

	return configuration
}
