package config

import (
	"errors"
	"log"
	"os"
)

type RegionalConfig struct {
	SKU          string
	Locale       string
	NvidiaLocale string
	Currency     string
	TestSKU      string
}

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type TwilioConfig struct {
	AccountSID        string
	Token             string
	SourceNumber      string
	DestinationNumber string
}

type DiscordConfig struct {
	WebhookURL string
}

type TelegramConfig struct {
	APIKey string
	ChatID string
}

type Config struct {
	Locale   string
	Currency string
	Delay    int64

	SKU            string
	TestSKU        string
	TwilioConfig   TwilioConfig
	TwitterConfig  TwitterConfig
	DiscordConfig  DiscordConfig
	TelegramConfig TelegramConfig
}

// Hardcoded SKU to locale/currency mappings to avoid user pain of having to lookup and enter these.
var regionalConfig = map[string]RegionalConfig{
	"AUT": {
		SKU:          "5440853700",
		Locale:       "de_de",
		NvidiaLocale: "de_de",
		Currency:     "EUR",
	},
	"BEL": {
		SKU:          "5438795700",
		Locale:       "fr_fr",
		NvidiaLocale: "fr_fr",
		Currency:     "EUR",
	},
	"CAN": {
		SKU:          "5438481700",
		Locale:       "en_us",
		NvidiaLocale: "en_ca",
		Currency:     "CAN",
	},
	"CZE": {
		SKU:          "5438793800",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"DNK": {
		SKU:          "5438793300",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"FIN": {
		SKU:          "5438793300",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"FRA": {
		SKU:          "5438795200",
		Locale:       "fr_fr",
		NvidiaLocale: "fr_fr",
		Currency:     "EUR",
	},
	"DEU": {
		SKU:          "5438792300",
		Locale:       "de_de",
		NvidiaLocale: "de_de",
		Currency:     "EUR",
	},
	"USA": {
		SKU:          "5438481700",
		Locale:       "en_us",
		NvidiaLocale: "en_us",
		Currency:     "USD",
		TestSKU:      "5379432500",
	},
	"GBR": {
		SKU:          "5438792800",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "GBP",
	},
	"IRL": {
		SKU:          "5438792800",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "GBP",
	},
	"ITA": {
		SKU:          "5438796200",
		Locale:       "it_it",
		NvidiaLocale: "it_it",
		Currency:     "EUR",
	},
	"SWE": {
		SKU:          "5438798100",
		Locale:       "sv_SE",
		NvidiaLocale: "sv_se",
		Currency:     "SEK",
	},
	"LUX": {
		SKU:          "5438795700",
		Locale:       "fr_fr",
		NvidiaLocale: "fr_fr",
		Currency:     "EUR",
	},
	"POL": {
		SKU:          "5438797700",
		Locale:       "pl_pl",
		NvidiaLocale: "pl_pl",
		Currency:     "PLN",
	},
	"PRT": {
		SKU:          "5438794300",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"ESP": {
		SKU:          "5438794800",
		Locale:       "es_es",
		NvidiaLocale: "es_es",
		Currency:     "EUR",
	},
	"NOR": {
		SKU:          "5438797200",
		Locale:       "no_NO",
		NvidiaLocale: "no_NO",
		Currency:     "NOK",
	},
	"NLD": {
		SKU:          "5438796700",
		Locale:       "nl_nl",
		NvidiaLocale: "nl_nl",
		Currency:     "EUR",
	},
}

//TwitterConfig Generates TwitterConfiguration for application from environmental variables.
func GetTwitterConfig() TwitterConfig {
	configuration := TwitterConfig{}

	consumerKey, consumerKeyOk := os.LookupEnv("TWITTER_CONSUMER_KEY")
	if consumerKeyOk == false {
		log.Fatal("TWITTER_CONSUMER_KEY Environment Variable is unset, exiting.")
	}

	configuration.ConsumerKey = consumerKey

	consumerSecret, consumerSecretOk := os.LookupEnv("TWITTER_CONSUMER_SECRET")
	if consumerSecretOk == false {
		log.Fatal("TWITTER_CONSUMER_SECRET Environment Variable is unset, exiting.")
	}

	configuration.ConsumerSecret = consumerSecret

	accessToken, accessTokenOk := os.LookupEnv("TWITTER_ACCESS_TOKEN")
	if accessTokenOk == false {
		log.Fatal("TWITTER_ACCESS_TOKEN Environment Variable is unset, exiting.")
	}

	configuration.AccessToken = accessToken

	accessSecret, accessSecretOk := os.LookupEnv("TWITTER_ACCESS_SECRET")
	if accessSecretOk == false {
		log.Fatal("TWILIO_DESTINATION_NUMBER Environment Variable is unset, exiting.")
	}

	configuration.AccessSecret = accessSecret

	return configuration
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

//GetDiscordConfig Generates DiscordConfiguration for application from environmental variables.
func GetDiscordConfig() DiscordConfig {
	configuration := DiscordConfig{}

	webhookURL, webhookURLOk := os.LookupEnv("DISCORD_WEBHOOK_URL")
	if webhookURLOk == false {
		log.Fatal("DISCORD_WEBHOOK_URL Environment Variable is unset, exiting.")
	}

	configuration.WebhookURL = webhookURL

	return configuration
}

//GetTelegramConfig Generates TelegramConfiguration for application from environmental variables.
func GetTelegramConfig() TelegramConfig {
	configuration := TelegramConfig{}

	APIKey, APIKeyOk := os.LookupEnv("TELEGRAM_API_KEY")
	if APIKeyOk == false {
		log.Fatal("TELEGRAM_API_KEY Environment Variable is unset, exiting.")
	}

	ChatID, ChatIDOk := os.LookupEnv("TELEGRAM_CHAT_ID")
	if ChatIDOk == false {
		log.Fatal("TELEGRAM_CHAT_ID Environment Variable is unset, exiting.")
	}

	configuration.APIKey = APIKey
	configuration.ChatID = ChatID

	return configuration
}

//GetConfig Generates Configuration for application from environmental variables.
func GetConfig(region string, delay int64, smsEnabled bool, discordEnabled bool, twitterEnabled bool, telegramEnabled bool) (*Config, error) {
	if regionConfig, ok := regionalConfig[region]; ok {
		configuration := Config{}

		configuration.SKU = regionConfig.SKU
		configuration.Delay = delay
		configuration.TestSKU = regionConfig.TestSKU
		configuration.Locale = regionConfig.Locale
		configuration.Currency = regionConfig.Currency

		if smsEnabled == true {
			configuration.TwilioConfig = GetTwilioConfig()
		}

		if discordEnabled == true {
			configuration.DiscordConfig = GetDiscordConfig()
		}

		if twitterEnabled == true {
			configuration.TwitterConfig = GetTwitterConfig()
		}

		if telegramEnabled == true {
			configuration.TelegramConfig = GetTelegramConfig()
		}

		return &configuration, nil
	}

	return nil, errors.New("unsupported region provided")
}
